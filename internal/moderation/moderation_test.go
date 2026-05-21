package moderation

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sumi-devs/canopy-social/canopy/internal/accounts"
	"github.com/sumi-devs/canopy-social/canopy/internal/auth"
)

type mockModerationRepository struct {
	reports  map[string]*Report
	blocks   map[string]*InstanceBlock
	settings map[string]*ServerSetting
	invites  map[string]*ServerInvite
}

func newMockModerationRepository() *mockModerationRepository {
	return &mockModerationRepository{
		reports:  make(map[string]*Report),
		blocks:   make(map[string]*InstanceBlock),
		settings: make(map[string]*ServerSetting),
		invites:  make(map[string]*ServerInvite),
	}
}

func (m *mockModerationRepository) CreateReport(ctx context.Context, report *Report) (*Report, error) {
	m.reports[report.ID] = report
	return report, nil
}

func (m *mockModerationRepository) GetReportByID(ctx context.Context, id string) (*Report, error) {
	rep, ok := m.reports[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return rep, nil
}

func (m *mockModerationRepository) ListReports(ctx context.Context, status string, limit, offset int) ([]*Report, error) {
	var res []*Report
	for _, r := range m.reports {
		if r.Status == status {
			res = append(res, r)
		}
	}
	return res, nil
}

func (m *mockModerationRepository) UpdateReport(ctx context.Context, report *Report) error {
	m.reports[report.ID] = report
	return nil
}

func (m *mockModerationRepository) CreateInstanceBlock(ctx context.Context, block *InstanceBlock) (*InstanceBlock, error) {
	m.blocks[block.ID] = block
	return block, nil
}

func (m *mockModerationRepository) GetInstanceBlockByID(ctx context.Context, id string) (*InstanceBlock, error) {
	b, ok := m.blocks[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return b, nil
}

func (m *mockModerationRepository) GetInstanceBlockByDomain(ctx context.Context, domain string) (*InstanceBlock, error) {
	for _, b := range m.blocks {
		if b.Domain == domain {
			return b, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockModerationRepository) ListInstanceBlocks(ctx context.Context, limit, offset int) ([]*InstanceBlock, error) {
	var res []*InstanceBlock
	for _, b := range m.blocks {
		res = append(res, b)
	}
	return res, nil
}

func (m *mockModerationRepository) DeleteInstanceBlock(ctx context.Context, id string) error {
	if _, ok := m.blocks[id]; !ok {
		return errors.New("not found")
	}
	delete(m.blocks, id)
	return nil
}

func (m *mockModerationRepository) GetSetting(ctx context.Context, key string) (*ServerSetting, error) {
	s, ok := m.settings[key]
	if !ok {
		return nil, errors.New("not found")
	}
	return s, nil
}

func (m *mockModerationRepository) PutSetting(ctx context.Context, key string, value json.RawMessage) error {
	m.settings[key] = &ServerSetting{
		Key:       key,
		Value:     value,
		UpdatedAt: time.Now(),
	}
	return nil
}

func (m *mockModerationRepository) CreateInvite(ctx context.Context, invite *ServerInvite) (*ServerInvite, error) {
	m.invites[invite.ID] = invite
	return invite, nil
}

func (m *mockModerationRepository) GetInviteByID(ctx context.Context, id string) (*ServerInvite, error) {
	i, ok := m.invites[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return i, nil
}

func (m *mockModerationRepository) GetInviteByToken(ctx context.Context, token string) (*ServerInvite, error) {
	for _, i := range m.invites {
		if i.Token == token {
			return i, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockModerationRepository) ListInvites(ctx context.Context, limit, offset int) ([]*ServerInvite, error) {
	var res []*ServerInvite
	for _, i := range m.invites {
		res = append(res, i)
	}
	return res, nil
}

func (m *mockModerationRepository) UpdateInvite(ctx context.Context, invite *ServerInvite) error {
	m.invites[invite.ID] = invite
	return nil
}

type mockAccountsRepository struct {
	accounts map[string]*accounts.Account
}

func (m *mockAccountsRepository) GetByID(ctx context.Context, id string) (*accounts.Account, error) {
	acc, ok := m.accounts[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return acc, nil
}
func (m *mockAccountsRepository) GetByURI(ctx context.Context, uri string) (*accounts.Account, error) {
	return nil, nil
}
func (m *mockAccountsRepository) GetByUsername(ctx context.Context, username string) (*accounts.Account, error) {
	for _, a := range m.accounts {
		if a.Username == username {
			return a, nil
		}
	}
	return nil, errors.New("not found")
}
func (m *mockAccountsRepository) GetByUsernameAndDomain(ctx context.Context, username, domain string) (*accounts.Account, error) {
	return nil, nil
}
func (m *mockAccountsRepository) GetByEmail(ctx context.Context, email string) (*accounts.Account, error) {
	return nil, nil
}
func (m *mockAccountsRepository) Create(ctx context.Context, account *accounts.Account) (*accounts.Account, error) {
	m.accounts[account.ID] = account
	return account, nil
}
func (m *mockAccountsRepository) UpdateProfile(ctx context.Context, id string, params *accounts.UpdateProfileParams) (*accounts.Account, error) {
	return nil, nil
}
func (m *mockAccountsRepository) ListLocal(ctx context.Context, limit, offset int) ([]*accounts.Account, error) {
	return nil, nil
}
func (m *mockAccountsRepository) SearchByUsername(ctx context.Context, query string, limit int) ([]*accounts.Account, error) {
	return nil, nil
}
func (m *mockAccountsRepository) Follow(ctx context.Context, followerID, followingID, status string) error {
	return nil
}
func (m *mockAccountsRepository) Unfollow(ctx context.Context, followerID, followingID string) error {
	return nil
}
func (m *mockAccountsRepository) AcceptFollow(ctx context.Context, followerID, followingID string) error {
	return nil
}
func (m *mockAccountsRepository) RejectFollow(ctx context.Context, followerID, followingID string) error {
	return nil
}
func (m *mockAccountsRepository) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	return false, nil
}
func (m *mockAccountsRepository) ListFollowers(ctx context.Context, accountID string, limit, offset int) ([]*accounts.Account, error) {
	return nil, nil
}
func (m *mockAccountsRepository) ListFollowing(ctx context.Context, accountID string, limit, offset int) ([]*accounts.Account, error) {
	return nil, nil
}
func (m *mockAccountsRepository) ListPendingRequests(ctx context.Context, accountID string, limit, offset int) ([]*accounts.Account, error) {
	return nil, nil
}
func (m *mockAccountsRepository) Block(ctx context.Context, accountID, targetID string) error {
	return nil
}
func (m *mockAccountsRepository) Unblock(ctx context.Context, accountID, targetID string) error {
	return nil
}
func (m *mockAccountsRepository) IsBlocking(ctx context.Context, accountID, targetID string) (bool, error) {
	return false, nil
}
func (m *mockAccountsRepository) Mute(ctx context.Context, accountID, targetID string, hideNotifications bool) error {
	return nil
}
func (m *mockAccountsRepository) Unmute(ctx context.Context, accountID, targetID string) error {
	return nil
}
func (m *mockAccountsRepository) IsMuting(ctx context.Context, accountID, targetID string) (bool, error) {
	return false, nil
}
func (m *mockAccountsRepository) ListBlocks(ctx context.Context, accountID string, limit, offset int) ([]*accounts.Account, error) {
	return nil, nil
}
func (m *mockAccountsRepository) ListMutes(ctx context.Context, accountID string, limit, offset int) ([]*accounts.Account, error) {
	return nil, nil
}
func (m *mockAccountsRepository) IncrementFollowersCount(ctx context.Context, id string) error {
	return nil
}
func (m *mockAccountsRepository) DecrementFollowersCount(ctx context.Context, id string) error {
	return nil
}
func (m *mockAccountsRepository) IncrementFollowingCount(ctx context.Context, id string) error {
	return nil
}
func (m *mockAccountsRepository) DecrementFollowingCount(ctx context.Context, id string) error {
	return nil
}
func (m *mockAccountsRepository) SetSuspended(ctx context.Context, id string, suspended bool) error {
	return nil
}
func (m *mockAccountsRepository) SetSilenced(ctx context.Context, id string, silenced bool) error {
	return nil
}
func (m *mockAccountsRepository) SetRole(ctx context.Context, id string, role string) error {
	return nil
}

func TestModeration_ServiceAndHandlers(t *testing.T) {
	repo := newMockModerationRepository()
	accsRepo := &mockAccountsRepository{
		accounts: map[string]*accounts.Account{
			"reporter1": {ID: "reporter1", Username: "rep1", IsLocal: true, Role: "user"},
			"target1":   {ID: "target1", Username: "targ1", IsLocal: true, Role: "user"},
			"admin1":    {ID: "admin1", Username: "adm1", IsLocal: true, Role: "admin"},
		},
	}
	accsSvc := accounts.NewService(accsRepo)
	svc := NewService(repo, accsSvc)
	h := NewHandler(svc)

	ctx := context.Background()
	targetID := "target1"
	category := "spam"
	comment := "this user is posting spam"
	report, err := svc.ReportContent(ctx, "reporter1", &targetID, nil, nil, category, comment)
	if err != nil {
		t.Fatalf("failed to report: %v", err)
	}
	if report.ReporterID != "reporter1" || *report.TargetAccountID != "target1" || report.Category != "spam" {
		t.Errorf("unexpected report details: %+v", report)
	}

	reports, err := svc.ListReports(ctx, "open", 10, 0)
	if err != nil {
		t.Fatalf("failed to list reports: %v", err)
	}
	if len(reports) != 1 {
		t.Errorf("expected 1 report, got %d", len(reports))
	}

	resolved, err := svc.ResolveReport(ctx, report.ID, "admin1", "suspended user")
	if err != nil {
		t.Fatalf("failed to resolve report: %v", err)
	}
	if resolved.Status != "resolved" || *resolved.ResolvedBy != "admin1" || *resolved.ActionTaken != "suspended user" {
		t.Errorf("unexpected resolved details: %+v", resolved)
	}

	block, err := svc.BlockInstance(ctx, "admin1", "spammy.domain", "suspend", "spamming", true, false)
	if err != nil {
		t.Fatalf("failed to block instance: %v", err)
	}
	if block.Domain != "spammy.domain" || block.Severity != "suspend" || *block.Reason != "spamming" {
		t.Errorf("unexpected block details: %+v", block)
	}

	blocks, err := svc.ListInstanceBlocks(ctx, 10, 0)
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 1 {
		t.Errorf("expected 1 block, got %d", len(blocks))
	}

	if err := svc.UnblockInstance(ctx, block.ID); err != nil {
		t.Fatalf("failed to unblock: %v", err)
	}
	if len(repo.blocks) != 0 {
		t.Errorf("expected 0 blocks, got %d", len(repo.blocks))
	}

	type ConfigData struct {
		SiteName string `json:"site_name"`
	}
	cfg := ConfigData{SiteName: "Canopy Test Site"}
	if err := svc.PutSetting(ctx, "site_config", cfg); err != nil {
		t.Fatalf("failed to put setting: %v", err)
	}
	var loaded ConfigData
	if err := svc.GetSetting(ctx, "site_config", &loaded); err != nil {
		t.Fatalf("failed to get setting: %v", err)
	}
	if loaded.SiteName != "Canopy Test Site" {
		t.Errorf("expected site name 'Canopy Test Site', got '%s'", loaded.SiteName)
	}

	maxUses := 3
	invite, err := svc.CreateInvite(ctx, "admin1", &maxUses, nil)
	if err != nil {
		t.Fatalf("failed to create invite: %v", err)
	}
	if invite.MaxUses == nil || *invite.MaxUses != 3 || invite.UsesCount != 0 {
		t.Errorf("unexpected invite: %+v", invite)
	}

	if err := svc.UseInvite(ctx, invite.Token); err != nil {
		t.Fatalf("failed to use invite: %v", err)
	}
	if invite.UsesCount != 1 {
		t.Errorf("expected uses count 1, got %d", invite.UsesCount)
	}

	if err := svc.RevokeInvite(ctx, invite.ID); err != nil {
		t.Fatalf("failed to revoke invite: %v", err)
	}
	if !invite.IsRevoked {
		t.Error("expected invite to be marked revoked")
	}

	if err := svc.UseInvite(ctx, invite.Token); err == nil {
		t.Error("expected error when using revoked invite")
	}

	jwtMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), auth.ContextKeyAccountID, "admin1")
			ctx = context.WithValue(ctx, auth.ContextKeyRole, "admin")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	router := chi.NewRouter()
	h.RegisterRoutes(router, jwtMiddleware)
	server := httptest.NewServer(router)
	defer server.Close()

	reportPayload := `{"account_id":"target1","category":"spam","comment":"spam comments"}`
	req, _ := http.NewRequest("POST", server.URL+"/api/v1/reports", strings.NewReader(reportPayload))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed HTTP POST report: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	req, _ = http.NewRequest("GET", server.URL+"/api/v1/admin/reports", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed HTTP GET reports: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	blockPayload := `{"domain":"malicious.com","severity":"suspend","reason":"phishing","reject_media":true,"reject_reports":true}`
	req, _ = http.NewRequest("POST", server.URL+"/api/v1/admin/instance_blocks", strings.NewReader(blockPayload))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed HTTP POST block: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}
}
