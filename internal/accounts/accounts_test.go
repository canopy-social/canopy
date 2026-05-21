package accounts

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sumi-devs/canopy-social/canopy/internal/auth"
)

type mockRepository struct {
	accounts map[string]*Account
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		accounts: make(map[string]*Account),
	}
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*Account, error) {
	acc, ok := m.accounts[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return acc, nil
}

func (m *mockRepository) GetByURI(ctx context.Context, uri string) (*Account, error) {
	return nil, nil
}

func (m *mockRepository) GetByUsername(ctx context.Context, username string) (*Account, error) {
	for _, a := range m.accounts {
		if a.Username == username {
			return a, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockRepository) GetByUsernameAndDomain(ctx context.Context, username, domain string) (*Account, error) {
	return nil, nil
}

func (m *mockRepository) GetByEmail(ctx context.Context, email string) (*Account, error) {
	return nil, nil
}

func (m *mockRepository) Create(ctx context.Context, account *Account) (*Account, error) {
	m.accounts[account.ID] = account
	return account, nil
}

func (m *mockRepository) UpdateProfile(ctx context.Context, id string, params *UpdateProfileParams) (*Account, error) {
	return nil, nil
}

func (m *mockRepository) ListLocal(ctx context.Context, limit, offset int) ([]*Account, error) {
	var res []*Account
	for _, a := range m.accounts {
		if a.IsLocal {
			res = append(res, a)
		}
	}
	return res, nil
}

func (m *mockRepository) SearchByUsername(ctx context.Context, query string, limit int) ([]*Account, error) {
	return nil, nil
}

func (m *mockRepository) Follow(ctx context.Context, followerID, followingID, status string) error {
	return nil
}

func (m *mockRepository) Unfollow(ctx context.Context, followerID, followingID string) error {
	return nil
}

func (m *mockRepository) AcceptFollow(ctx context.Context, followerID, followingID string) error {
	return nil
}

func (m *mockRepository) RejectFollow(ctx context.Context, followerID, followingID string) error {
	return nil
}

func (m *mockRepository) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	return false, nil
}

func (m *mockRepository) ListFollowers(ctx context.Context, accountID string, limit, offset int) ([]*Account, error) {
	return nil, nil
}

func (m *mockRepository) ListFollowing(ctx context.Context, accountID string, limit, offset int) ([]*Account, error) {
	return nil, nil
}

func (m *mockRepository) ListPendingRequests(ctx context.Context, accountID string, limit, offset int) ([]*Account, error) {
	return nil, nil
}

func (m *mockRepository) Block(ctx context.Context, accountID, targetID string) error {
	return nil
}

func (m *mockRepository) Unblock(ctx context.Context, accountID, targetID string) error {
	return nil
}

func (m *mockRepository) IsBlocking(ctx context.Context, accountID, targetID string) (bool, error) {
	return false, nil
}

func (m *mockRepository) Mute(ctx context.Context, accountID, targetID string, hideNotifications bool) error {
	return nil
}

func (m *mockRepository) Unmute(ctx context.Context, accountID, targetID string) error {
	return nil
}

func (m *mockRepository) IsMuting(ctx context.Context, accountID, targetID string) (bool, error) {
	return false, nil
}

func (m *mockRepository) ListBlocks(ctx context.Context, accountID string, limit, offset int) ([]*Account, error) {
	return nil, nil
}

func (m *mockRepository) ListMutes(ctx context.Context, accountID string, limit, offset int) ([]*Account, error) {
	return nil, nil
}

func (m *mockRepository) IncrementFollowersCount(ctx context.Context, id string) error {
	return nil
}

func (m *mockRepository) DecrementFollowersCount(ctx context.Context, id string) error {
	return nil
}

func (m *mockRepository) IncrementFollowingCount(ctx context.Context, id string) error {
	return nil
}

func (m *mockRepository) DecrementFollowingCount(ctx context.Context, id string) error {
	return nil
}

func (m *mockRepository) SetSuspended(ctx context.Context, id string, suspended bool) error {
	acc, ok := m.accounts[id]
	if !ok {
		return errors.New("not found")
	}
	acc.IsSuspended = suspended
	return nil
}

func (m *mockRepository) SetSilenced(ctx context.Context, id string, silenced bool) error {
	acc, ok := m.accounts[id]
	if !ok {
		return errors.New("not found")
	}
	acc.IsSilenced = silenced
	return nil
}

func (m *mockRepository) SetRole(ctx context.Context, id string, role string) error {
	acc, ok := m.accounts[id]
	if !ok {
		return errors.New("not found")
	}
	acc.Role = role
	return nil
}

func TestAdmin_ServiceActions(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)
	ctx := context.Background()

	acc := &Account{ID: "test-user", Username: "test", Role: "user", IsLocal: true}
	repo.accounts[acc.ID] = acc

	if err := svc.SetSuspended(ctx, "test-user", true); err != nil {
		t.Fatalf("failed to suspend: %v", err)
	}
	if !acc.IsSuspended {
		t.Error("expected account to be suspended")
	}

	if err := svc.SetSilenced(ctx, "test-user", true); err != nil {
		t.Fatalf("failed to silence: %v", err)
	}
	if !acc.IsSilenced {
		t.Error("expected account to be silenced")
	}

	if err := svc.SetRole(ctx, "test-user", "moderator"); err != nil {
		t.Fatalf("failed to set role: %v", err)
	}
	if acc.Role != "moderator" {
		t.Errorf("expected role 'moderator', got '%s'", acc.Role)
	}

	if err := svc.SetRole(ctx, "test-user", "invalid"); err == nil {
		t.Error("expected error for invalid role")
	}
}

func TestAdmin_HandlerActions(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)
	h := NewHandler(svc)

	acc1 := &Account{ID: "u1", Username: "user1", Role: "user", IsLocal: true}
	acc2 := &Account{ID: "u2", Username: "user2", Role: "user", IsLocal: true}
	repo.accounts[acc1.ID] = acc1
	repo.accounts[acc2.ID] = acc2

	jwtMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), auth.ContextKeyAccountID, "admin-user")
			ctx = context.WithValue(ctx, auth.ContextKeyRole, "admin")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	r := chi.NewRouter()
	h.RegisterRoutes(r, jwtMiddleware)
	server := httptest.NewServer(r)
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL+"/api/v1/admin/accounts", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to GET accounts: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var list []*Account
	json.NewDecoder(resp.Body).Decode(&list)
	if len(list) != 2 {
		t.Errorf("expected 2 local accounts, got %d", len(list))
	}

	payload := `{"type":"suspend"}`
	req, _ = http.NewRequest("POST", server.URL+"/api/v1/admin/accounts/u1/action", strings.NewReader(payload))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed action POST: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var updated Account
	json.NewDecoder(resp.Body).Decode(&updated)
	if !acc1.IsSuspended {
		t.Error("expected account to be marked suspended in repo")
	}

	payload = `{"type":"change_role","role":"moderator"}`
	req, _ = http.NewRequest("POST", server.URL+"/api/v1/admin/accounts/u2/action", strings.NewReader(payload))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed action POST: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var updated2 Account
	json.NewDecoder(resp.Body).Decode(&updated2)
	if updated2.Role != "moderator" {
		t.Errorf("expected moderator role, got %s", updated2.Role)
	}
}
