package moderation

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/sumi-devs/canopy-social/canopy/internal/accounts"
	"github.com/sumi-devs/canopy-social/canopy/pkg/ulid"
)

type Service struct {
	repo     Repository
	accounts *accounts.Service
}

func NewService(repo Repository, accounts *accounts.Service) *Service {
	return &Service{
		repo:     repo,
		accounts: accounts,
	}
}

func (s *Service) ReportContent(ctx context.Context, reporterID string, targetAccountID, targetPostID, targetEssayID *string, category, comment string) (*Report, error) {
	_, err := s.accounts.GetByID(ctx, reporterID)
	if err != nil {
		return nil, fmt.Errorf("reporter not found: %w", err)
	}

	if targetAccountID == nil && targetPostID == nil && targetEssayID == nil {
		return nil, errors.New("at least one report target is required")
	}

	if targetAccountID != nil {
		_, err := s.accounts.GetByID(ctx, *targetAccountID)
		if err != nil {
			return nil, fmt.Errorf("target account not found: %w", err)
		}
	}

	var commentPtr *string
	if comment != "" {
		commentPtr = &comment
	}

	report := &Report{
		ID:              ulid.New(),
		ReporterID:      reporterID,
		TargetAccountID: targetAccountID,
		TargetPostID:    targetPostID,
		TargetEssayID:   targetEssayID,
		Category:        category,
		Comment:         commentPtr,
		Status:          "open",
		CreatedAt:       time.Now(),
	}

	return s.repo.CreateReport(ctx, report)
}

func (s *Service) ListReports(ctx context.Context, status string, limit, offset int) ([]*Report, error) {
	reports, err := s.repo.ListReports(ctx, status, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, r := range reports {
		if reporter, err := s.accounts.GetByID(ctx, r.ReporterID); err == nil {
			r.Reporter = reporter
		}
		if r.TargetAccountID != nil {
			if target, err := s.accounts.GetByID(ctx, *r.TargetAccountID); err == nil {
				r.TargetAccount = target
			}
		}
	}

	return reports, nil
}

func (s *Service) ResolveReport(ctx context.Context, id, resolverID, actionTaken string) (*Report, error) {
	report, err := s.repo.GetReportByID(ctx, id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	report.Status = "resolved"
	report.ResolvedBy = &resolverID
	report.ResolvedAt = &now
	report.ActionTaken = &actionTaken

	if err := s.repo.UpdateReport(ctx, report); err != nil {
		return nil, err
	}

	return report, nil
}

func (s *Service) BlockInstance(ctx context.Context, creatorID, domain, severity, reason string, rejectMedia, rejectReports bool) (*InstanceBlock, error) {
	if domain == "" {
		return nil, errors.New("domain is required")
	}

	if severity != "silence" && severity != "suspend" {
		return nil, errors.New("invalid block severity")
	}

	var reasonPtr *string
	if reason != "" {
		reasonPtr = &reason
	}

	block := &InstanceBlock{
		ID:            ulid.New(),
		Domain:        domain,
		Severity:      severity,
		Reason:        reasonPtr,
		RejectMedia:   rejectMedia,
		RejectReports: rejectReports,
		CreatedBy:     &creatorID,
		CreatedAt:     time.Now(),
	}

	return s.repo.CreateInstanceBlock(ctx, block)
}

func (s *Service) ListInstanceBlocks(ctx context.Context, limit, offset int) ([]*InstanceBlock, error) {
	return s.repo.ListInstanceBlocks(ctx, limit, offset)
}

func (s *Service) UnblockInstance(ctx context.Context, id string) error {
	return s.repo.DeleteInstanceBlock(ctx, id)
}

func (s *Service) GetSetting(ctx context.Context, key string, target interface{}) error {
	setting, err := s.repo.GetSetting(ctx, key)
	if err != nil {
		return err
	}

	return json.Unmarshal(setting.Value, target)
}

func (s *Service) PutSetting(ctx context.Context, key string, value interface{}) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.repo.PutSetting(ctx, key, json.RawMessage(raw))
}

func (s *Service) CreateInvite(ctx context.Context, creatorID string, maxUses *int, expiresAt *time.Time) (*ServerInvite, error) {
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(tokenBytes)

	invite := &ServerInvite{
		ID:        ulid.New(),
		CreatedBy: creatorID,
		Token:     token,
		MaxUses:   maxUses,
		UsesCount: 0,
		ExpiresAt: expiresAt,
		IsRevoked: false,
		CreatedAt: time.Now(),
	}

	return s.repo.CreateInvite(ctx, invite)
}

func (s *Service) ListInvites(ctx context.Context, limit, offset int) ([]*ServerInvite, error) {
	return s.repo.ListInvites(ctx, limit, offset)
}

func (s *Service) RevokeInvite(ctx context.Context, id string) error {
	invite, err := s.repo.GetInviteByID(ctx, id)
	if err != nil {
		return err
	}

	invite.IsRevoked = true
	return s.repo.UpdateInvite(ctx, invite)
}

func (s *Service) UseInvite(ctx context.Context, token string) error {
	invite, err := s.repo.GetInviteByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invite not found: %w", err)
	}

	if invite.IsRevoked {
		return errors.New("invite has been revoked")
	}

	if invite.ExpiresAt != nil && invite.ExpiresAt.Before(time.Now()) {
		return errors.New("invite has expired")
	}

	if invite.MaxUses != nil && invite.UsesCount >= *invite.MaxUses {
		return errors.New("invite maximum uses reached")
	}

	invite.UsesCount++
	return s.repo.UpdateInvite(ctx, invite)
}
