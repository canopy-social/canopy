package timeline

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sumi-devs/canopy-social/canopy/internal/accounts"
	"github.com/sumi-devs/canopy-social/canopy/internal/posts"
	"github.com/sumi-devs/canopy-social/canopy/pkg/config"
)

type mockPostRepo struct {
	posts map[string]*posts.Post
}

func (m *mockPostRepo) GetByID(ctx context.Context, id string) (*posts.Post, error) {
	p, ok := m.posts[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return p, nil
}

func (m *mockPostRepo) ListPostsByAccountWithBoosts(ctx context.Context, accountID string, limit, offset int) ([]*posts.Post, error) {
	var res []*posts.Post
	for _, p := range m.posts {
		if p.AccountID == accountID {
			res = append(res, p)
		}
	}
	sortPostsDesc(res)
	if len(res) > limit {
		res = res[:limit]
	}
	return res, nil
}

func (m *mockPostRepo) ListPublicTimeline(ctx context.Context, local bool, limit, offset int) ([]*posts.Post, error) {
	var res []*posts.Post
	for _, p := range m.posts {
		if p.Visibility == "public" && p.ReplyToID == nil && p.BoostOfID == nil {
			if !local || p.IsLocal {
				res = append(res, p)
			}
		}
	}
	sortPostsDesc(res)
	if len(res) > limit {
		res = res[:limit]
	}
	return res, nil
}

type mockAccountRepo struct {
	accounts  map[string]*accounts.Account
	followers map[string][]string
	following map[string][]string
}

func (m *mockAccountRepo) GetByID(ctx context.Context, id string) (*accounts.Account, error) {
	a, ok := m.accounts[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return a, nil
}

func (m *mockAccountRepo) ListFollowers(ctx context.Context, accountID string, limit, offset int) ([]*accounts.Account, error) {
	ids := m.followers[accountID]
	var res []*accounts.Account
	for _, id := range ids {
		if a, ok := m.accounts[id]; ok {
			res = append(res, a)
		}
	}
	return res, nil
}

func (m *mockAccountRepo) ListFollowing(ctx context.Context, accountID string, limit, offset int) ([]*accounts.Account, error) {
	ids := m.following[accountID]
	var res []*accounts.Account
	for _, id := range ids {
		if a, ok := m.accounts[id]; ok {
			res = append(res, a)
		}
	}
	return res, nil
}

func sortPostsDesc(p []*posts.Post) {
	sort.Slice(p, func(i, j int) bool {
		return p[i].CreatedAt.After(p[j].CreatedAt)
	})
}

func TestFanoutPost_Normal(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	rclient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	cfg := &config.Config{}

	ar := &mockAccountRepo{
		accounts: map[string]*accounts.Account{
			"author": {ID: "author", FollowersCount: 2},
			"f1":     {ID: "f1"},
			"f2":     {ID: "f2"},
		},
		followers: map[string][]string{
			"author": {"f1", "f2"},
		},
	}
	pr := &mockPostRepo{
		posts: make(map[string]*posts.Post),
	}

	svc := NewService(rclient, cfg, pr, ar, 5)

	post := &posts.Post{
		ID:        "post123",
		AccountID: "author",
		CreatedAt: time.Now(),
	}

	err = svc.FanoutPost(context.Background(), post)
	if err != nil {
		t.Fatalf("failed to fanout: %v", err)
	}

	for _, id := range []string{"author", "f1", "f2"} {
		exists := mr.Exists("timeline:home:" + id)
		if !exists {
			t.Errorf("expected timeline:home:%s to exist", id)
		}
		members, err := rclient.ZRange(context.Background(), "timeline:home:"+id, 0, -1).Result()
		if err != nil || len(members) != 1 || members[0] != "post123" {
			t.Errorf("invalid sorted set for %s", id)
		}
	}
}

func TestFanoutPost_Celebrity(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	rclient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	cfg := &config.Config{}

	ar := &mockAccountRepo{
		accounts: map[string]*accounts.Account{
			"author": {ID: "author", FollowersCount: 10},
			"f1":     {ID: "f1"},
			"f2":     {ID: "f2"},
		},
		followers: map[string][]string{
			"author": {"f1", "f2"},
		},
	}
	pr := &mockPostRepo{
		posts: make(map[string]*posts.Post),
	}

	svc := NewService(rclient, cfg, pr, ar, 5)

	post := &posts.Post{
		ID:        "post123",
		AccountID: "author",
		CreatedAt: time.Now(),
	}

	err = svc.FanoutPost(context.Background(), post)
	if err != nil {
		t.Fatalf("failed to fanout: %v", err)
	}

	if !mr.Exists("timeline:home:author") {
		t.Error("expected author's own timeline to exist")
	}
	if mr.Exists("timeline:home:f1") || mr.Exists("timeline:home:f2") {
		t.Error("expected followers' timelines NOT to have the celebrity post on write")
	}
}

func TestGetHomeTimeline(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	rclient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	cfg := &config.Config{}

	now := time.Now()
	post1 := &posts.Post{ID: "p1", AccountID: "user1", CreatedAt: now.Add(-5 * time.Minute)}
	post2 := &posts.Post{ID: "p2", AccountID: "celebrity", CreatedAt: now.Add(-2 * time.Minute)}
	post3 := &posts.Post{ID: "p3", AccountID: "user1", CreatedAt: now.Add(-10 * time.Minute)}

	ar := &mockAccountRepo{
		accounts: map[string]*accounts.Account{
			"me":        {ID: "me"},
			"user1":     {ID: "user1", FollowersCount: 1},
			"celebrity": {ID: "celebrity", FollowersCount: 10},
		},
		following: map[string][]string{
			"me": {"user1", "celebrity"},
		},
	}
	pr := &mockPostRepo{
		posts: map[string]*posts.Post{
			"p1": post1,
			"p2": post2,
			"p3": post3,
		},
	}

	svc := NewService(rclient, cfg, pr, ar, 5)

	rclient.ZAdd(context.Background(), "timeline:home:me", redis.Z{
		Score:  float64(post1.CreatedAt.UnixMilli()),
		Member: "p1",
	})
	rclient.ZAdd(context.Background(), "timeline:home:me", redis.Z{
		Score:  float64(post3.CreatedAt.UnixMilli()),
		Member: "p3",
	})

	resp, err := svc.GetHomeTimeline(context.Background(), "me", 10, "", "")
	if err != nil {
		t.Fatalf("failed to get home timeline: %v", err)
	}

	if len(resp.Data) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(resp.Data))
	}

	if resp.Data[0].ID != "p2" || resp.Data[1].ID != "p1" || resp.Data[2].ID != "p3" {
		t.Errorf("posts are not sorted correctly by CreatedAt descending")
	}

	if resp.PrevCursor != "p2" || resp.NextCursor != "p3" {
		t.Errorf("cursors are not set correctly: prev=%s, next=%s", resp.PrevCursor, resp.NextCursor)
	}
}

func TestGetHomeTimeline_CursorPagination(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	rclient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	cfg := &config.Config{}

	now := time.Now()
	p1 := &posts.Post{ID: "p1", AccountID: "u1", CreatedAt: now}
	p2 := &posts.Post{ID: "p2", AccountID: "u1", CreatedAt: now.Add(-1 * time.Minute)}
	p3 := &posts.Post{ID: "p3", AccountID: "u1", CreatedAt: now.Add(-2 * time.Minute)}

	ar := &mockAccountRepo{
		accounts: map[string]*accounts.Account{
			"me": {ID: "me"},
			"u1": {ID: "u1", FollowersCount: 1},
		},
		following: map[string][]string{
			"me": {"u1"},
		},
	}
	pr := &mockPostRepo{
		posts: map[string]*posts.Post{
			"p1": p1,
			"p2": p2,
			"p3": p3,
		},
	}

	svc := NewService(rclient, cfg, pr, ar, 5)

	rclient.ZAdd(context.Background(), "timeline:home:me", redis.Z{Score: float64(p1.CreatedAt.UnixMilli()), Member: "p1"})
	rclient.ZAdd(context.Background(), "timeline:home:me", redis.Z{Score: float64(p2.CreatedAt.UnixMilli()), Member: "p2"})
	rclient.ZAdd(context.Background(), "timeline:home:me", redis.Z{Score: float64(p3.CreatedAt.UnixMilli()), Member: "p3"})

	resp, err := svc.GetHomeTimeline(context.Background(), "me", 1, "p1", "")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if len(resp.Data) != 1 || resp.Data[0].ID != "p2" {
		t.Errorf("expected post p2, got %v", resp.Data)
	}

	resp, err = svc.GetHomeTimeline(context.Background(), "me", 10, "", "p3")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if len(resp.Data) != 2 || resp.Data[0].ID != "p1" || resp.Data[1].ID != "p2" {
		t.Errorf("expected posts p1 and p2, got %v", resp.Data)
	}
}

func TestGetPublicTimeline_Cached(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	rclient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	cfg := &config.Config{}

	p1 := &posts.Post{ID: "p1", Visibility: "public", CreatedAt: time.Now()}

	ar := &mockAccountRepo{}
	pr := &mockPostRepo{
		posts: map[string]*posts.Post{
			"p1": p1,
		},
	}

	svc := NewService(rclient, cfg, pr, ar, 5)

	res, err := svc.GetPublicTimeline(context.Background(), false, 20, 0)
	if err != nil || len(res) != 1 || res[0].ID != "p1" {
		t.Fatalf("failed initial call: %v", err)
	}

	cacheKey := "cache:public_timeline:false:20:0"
	if !mr.Exists(cacheKey) {
		t.Error("expected public timeline cache key to exist")
	}

	delete(pr.posts, "p1")

	res2, err := svc.GetPublicTimeline(context.Background(), false, 20, 0)
	if err != nil || len(res2) != 1 || res2[0].ID != "p1" {
		t.Errorf("expected post to be loaded from cache, got %v", res2)
	}
}

func TestGetPublicTimeline_Local(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	rclient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	cfg := &config.Config{}

	p1 := &posts.Post{ID: "p1", Visibility: "public", IsLocal: true, CreatedAt: time.Now()}
	p2 := &posts.Post{ID: "p2", Visibility: "public", IsLocal: false, CreatedAt: time.Now().Add(-1 * time.Minute)}

	ar := &mockAccountRepo{}
	pr := &mockPostRepo{
		posts: map[string]*posts.Post{
			"p1": p1,
			"p2": p2,
		},
	}

	svc := NewService(rclient, cfg, pr, ar, 5)

	res, err := svc.GetPublicTimeline(context.Background(), true, 20, 0)
	if err != nil || len(res) != 1 || res[0].ID != "p1" {
		t.Fatalf("failed local timeline call: %v", err)
	}

	res2, err := svc.GetPublicTimeline(context.Background(), false, 20, 0)
	if err != nil || len(res2) != 2 {
		t.Fatalf("failed federated timeline call: %v", err)
	}
}

func TestHandler_Routes(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	rclient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	cfg := &config.Config{}

	p1 := &posts.Post{ID: "p1", Visibility: "public", CreatedAt: time.Now()}
	ar := &mockAccountRepo{
		accounts: map[string]*accounts.Account{
			"user": {ID: "user"},
		},
	}
	pr := &mockPostRepo{
		posts: map[string]*posts.Post{
			"p1": p1,
		},
	}

	svc := NewService(rclient, cfg, pr, ar, 5)
	h := NewHandler(svc)

	jwtMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "account_id", "user")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	optionalJWT := func(next http.Handler) http.Handler {
		return next
	}

	r := chi.NewRouter()
	h.RegisterRoutes(r, jwtMiddleware, optionalJWT)

	server := httptest.NewServer(r)
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/timelines/public")
	if err != nil {
		t.Fatalf("HTTP GET failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}

	req, _ := http.NewRequest("GET", server.URL+"/api/v1/timelines/home", nil)
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("HTTP GET failed: %v", err)
	}
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp2.StatusCode)
	}

	var timelineResp TimelineResponse
	json.NewDecoder(resp2.Body).Decode(&timelineResp)
	if len(timelineResp.Data) != 0 {
		t.Errorf("expected empty timeline response, got %v", timelineResp.Data)
	}
}
