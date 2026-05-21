package notifications

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sumi-devs/canopy-social/canopy/internal/accounts"
	"github.com/sumi-devs/canopy-social/canopy/internal/auth"
	"github.com/sumi-devs/canopy-social/canopy/internal/posts"
)

type mockRepository struct {
	items map[string]*Notification
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		items: make(map[string]*Notification),
	}
}

func (m *mockRepository) Create(ctx context.Context, n *Notification) (*Notification, error) {
	m.items[n.ID] = n
	return n, nil
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*Notification, error) {
	n, ok := m.items[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return n, nil
}

func (m *mockRepository) List(ctx context.Context, accountID string, limit int, maxID, sinceID string, types []string, excludeTypes []string) ([]*Notification, error) {
	var res []*Notification
	for _, n := range m.items {
		if n.AccountID != accountID {
			continue
		}
		if n.DismissedAt != nil {
			continue
		}
		if len(types) > 0 {
			match := false
			for _, t := range types {
				if n.Type == t {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		if len(excludeTypes) > 0 {
			exclude := false
			for _, t := range excludeTypes {
				if n.Type == t {
					exclude = true
					break
				}
			}
			if exclude {
				continue
			}
		}
		if maxID != "" && n.ID >= maxID {
			continue
		}
		if sinceID != "" && n.ID <= sinceID {
			continue
		}
		res = append(res, n)
	}
	return res, nil
}

func (m *mockRepository) MarkRead(ctx context.Context, accountID string, id string) (*Notification, error) {
	n, ok := m.items[id]
	if !ok || n.AccountID != accountID {
		return nil, errors.New("not found")
	}
	now := time.Now()
	n.ReadAt = &now
	return n, nil
}

func (m *mockRepository) MarkAllRead(ctx context.Context, accountID string) error {
	now := time.Now()
	for _, n := range m.items {
		if n.AccountID == accountID {
			n.ReadAt = &now
		}
	}
	return nil
}

func (m *mockRepository) Dismiss(ctx context.Context, accountID string, id string) error {
	n, ok := m.items[id]
	if !ok || n.AccountID != accountID {
		return errors.New("not found")
	}
	now := time.Now()
	n.DismissedAt = &now
	return nil
}

type mockAccountLoader struct {
	accounts map[string]*accounts.Account
}

func (m *mockAccountLoader) GetByID(ctx context.Context, id string) (*accounts.Account, error) {
	acc, ok := m.accounts[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return acc, nil
}

func (m *mockAccountLoader) GetByUsername(ctx context.Context, username string) (*accounts.Account, error) {
	for _, acc := range m.accounts {
		if acc.Username == username {
			return acc, nil
		}
	}
	return nil, errors.New("not found")
}

type mockPostLoader struct {
	posts map[string]*posts.Post
}

func (m *mockPostLoader) GetByID(ctx context.Context, id string) (*posts.Post, error) {
	p, ok := m.posts[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return p, nil
}

func TestNotification_OnFollow(t *testing.T) {
	repo := newMockRepository()
	al := &mockAccountLoader{}
	pl := &mockPostLoader{}
	svc := NewService(repo, al, pl)
	ctx := context.Background()

	svc.OnFollow(ctx, "follower1", "following1", "accepted")

	if len(repo.items) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(repo.items))
	}
	var n *Notification
	for _, v := range repo.items {
		n = v
	}
	if n.Type != "follow" {
		t.Errorf("expected follow notification, got %s", n.Type)
	}
	if n.AccountID != "following1" {
		t.Errorf("expected account ID following1, got %s", n.AccountID)
	}
	if *n.FromAccountID != "follower1" {
		t.Errorf("expected from account follower1, got %s", *n.FromAccountID)
	}
}

func TestNotification_OnPostLiked(t *testing.T) {
	repo := newMockRepository()
	al := &mockAccountLoader{}
	pl := &mockPostLoader{
		posts: map[string]*posts.Post{
			"post1": {ID: "post1", AccountID: "author1"},
		},
	}
	svc := NewService(repo, al, pl)
	ctx := context.Background()

	svc.OnPostLiked(ctx, "post1", "liker1")

	if len(repo.items) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(repo.items))
	}
	var n *Notification
	for _, v := range repo.items {
		n = v
	}
	if n.Type != "favourite" {
		t.Errorf("expected favourite notification, got %s", n.Type)
	}
	if n.AccountID != "author1" {
		t.Errorf("expected recipient author1, got %s", n.AccountID)
	}
	if *n.PostID != "post1" {
		t.Errorf("expected post post1, got %s", *n.PostID)
	}
	if *n.FromAccountID != "liker1" {
		t.Errorf("expected from liker1, got %s", *n.FromAccountID)
	}
}

func TestNotification_OnPostBoosted(t *testing.T) {
	repo := newMockRepository()
	al := &mockAccountLoader{}
	pl := &mockPostLoader{
		posts: map[string]*posts.Post{
			"post1": {ID: "post1", AccountID: "author1"},
		},
	}
	svc := NewService(repo, al, pl)
	ctx := context.Background()

	svc.OnPostBoosted(ctx, "post1", "booster1")

	if len(repo.items) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(repo.items))
	}
	var n *Notification
	for _, v := range repo.items {
		n = v
	}
	if n.Type != "reblog" {
		t.Errorf("expected reblog notification, got %s", n.Type)
	}
}

func TestNotification_OnPostCreated_Reply(t *testing.T) {
	repo := newMockRepository()
	al := &mockAccountLoader{}
	parentID := "parentPost"
	pl := &mockPostLoader{
		posts: map[string]*posts.Post{
			parentID: {ID: parentID, AccountID: "parentAuthor"},
		},
	}
	svc := NewService(repo, al, pl)
	ctx := context.Background()

	post := &posts.Post{
		ID:        "replyPost",
		AccountID: "replyAuthor",
		ReplyToID: &parentID,
	}

	svc.OnPostCreated(ctx, post)

	if len(repo.items) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(repo.items))
	}
	var n *Notification
	for _, v := range repo.items {
		n = v
	}
	if n.Type != "mention" {
		t.Errorf("expected mention notification for reply, got %s", n.Type)
	}
	if n.AccountID != "parentAuthor" {
		t.Errorf("expected parentAuthor recipient, got %s", n.AccountID)
	}
}

func TestNotification_OnPostCreated_Mention(t *testing.T) {
	repo := newMockRepository()
	al := &mockAccountLoader{
		accounts: map[string]*accounts.Account{
			"alice": {ID: "aliceID", Username: "alice"},
		},
	}
	pl := &mockPostLoader{}
	svc := NewService(repo, al, pl)
	ctx := context.Background()

	post := &posts.Post{
		ID:          "post1",
		AccountID:   "authorID",
		ContentText: "hello @alice how are you?",
	}

	svc.OnPostCreated(ctx, post)

	if len(repo.items) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(repo.items))
	}
	var n *Notification
	for _, v := range repo.items {
		n = v
	}
	if n.Type != "mention" {
		t.Errorf("expected mention notification, got %s", n.Type)
	}
	if n.AccountID != "aliceID" {
		t.Errorf("expected aliceID recipient, got %s", n.AccountID)
	}
}

func TestNotification_Handler(t *testing.T) {
	repo := newMockRepository()
	al := &mockAccountLoader{
		accounts: map[string]*accounts.Account{
			"fromUser": {ID: "fromUser", Username: "fromUser"},
		},
	}
	pl := &mockPostLoader{
		posts: map[string]*posts.Post{
			"p1": {ID: "p1", AccountID: "user"},
		},
	}
	svc := NewService(repo, al, pl)
	h := NewHandler(svc)

	jwtMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), auth.ContextKeyAccountID, "user")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	r := chi.NewRouter()
	h.RegisterRoutes(r, jwtMiddleware)

	server := httptest.NewServer(r)
	defer server.Close()

	nID := "notif1"
	fromID := "fromUser"
	postID := "p1"
	repo.items[nID] = &Notification{
		ID:            nID,
		AccountID:     "user",
		Type:          "favourite",
		FromAccountID: &fromID,
		PostID:        &postID,
		CreatedAt:     time.Now(),
	}

	req, _ := http.NewRequest("GET", server.URL+"/api/v1/notifications", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed GET request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
	var list []*Notification
	json.NewDecoder(resp.Body).Decode(&list)
	if len(list) != 1 {
		t.Errorf("expected 1 notification in list, got %d", len(list))
	}
	if list[0].Account == nil || list[0].Account.Username != "fromUser" {
		t.Error("expected loaded sender account details")
	}
	if list[0].Status == nil || list[0].Status.ID != "p1" {
		t.Error("expected loaded status details")
	}

	req, _ = http.NewRequest("POST", server.URL+"/api/v1/notifications/notif1/read", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed POST read request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for read, got %d", resp.StatusCode)
	}
	var readN Notification
	json.NewDecoder(resp.Body).Decode(&readN)
	if readN.ReadAt == nil {
		t.Error("expected read_at to be non-nil")
	}

	req, _ = http.NewRequest("DELETE", server.URL+"/api/v1/notifications/notif1", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed DELETE dismiss request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for delete, got %d", resp.StatusCode)
	}
}
