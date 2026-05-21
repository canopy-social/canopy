package moderation

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sumi-devs/canopy-social/canopy/internal/accounts"
)

type Report struct {
	ID              string     `json:"id"`
	ReporterID      string     `json:"reporter_id"`
	TargetAccountID *string    `json:"target_account_id,omitempty"`
	TargetPostID    *string    `json:"target_post_id,omitempty"`
	TargetEssayID   *string    `json:"target_essay_id,omitempty"`
	Category        string     `json:"category"`
	Comment         *string    `json:"comment,omitempty"`
	Status          string     `json:"status"`
	ResolvedBy      *string    `json:"resolved_by,omitempty"`
	ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
	ActionTaken     *string    `json:"action_taken,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`

	Reporter      *accounts.Account `json:"reporter,omitempty"`
	TargetAccount *accounts.Account `json:"target_account,omitempty"`
}

type InstanceBlock struct {
	ID            string    `json:"id"`
	Domain        string    `json:"domain"`
	Severity      string    `json:"severity"`
	Reason        *string   `json:"reason,omitempty"`
	RejectMedia   bool      `json:"reject_media"`
	RejectReports bool      `json:"reject_reports"`
	CreatedBy     *string   `json:"created_by,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type ServerSetting struct {
	Key       string          `json:"key"`
	Value     json.RawMessage `json:"value"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ServerInvite struct {
	ID        string     `json:"id"`
	CreatedBy string     `json:"created_by"`
	Token     string     `json:"token"`
	MaxUses   *int       `json:"max_uses,omitempty"`
	UsesCount int        `json:"uses_count"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	IsRevoked bool       `json:"is_revoked"`
	CreatedAt time.Time  `json:"created_at"`
}

type Repository interface {
	CreateReport(ctx context.Context, report *Report) (*Report, error)
	GetReportByID(ctx context.Context, id string) (*Report, error)
	ListReports(ctx context.Context, status string, limit, offset int) ([]*Report, error)
	UpdateReport(ctx context.Context, report *Report) error

	CreateInstanceBlock(ctx context.Context, block *InstanceBlock) (*InstanceBlock, error)
	GetInstanceBlockByID(ctx context.Context, id string) (*InstanceBlock, error)
	GetInstanceBlockByDomain(ctx context.Context, domain string) (*InstanceBlock, error)
	ListInstanceBlocks(ctx context.Context, limit, offset int) ([]*InstanceBlock, error)
	DeleteInstanceBlock(ctx context.Context, id string) error

	GetSetting(ctx context.Context, key string) (*ServerSetting, error)
	PutSetting(ctx context.Context, key string, value json.RawMessage) error

	CreateInvite(ctx context.Context, invite *ServerInvite) (*ServerInvite, error)
	GetInviteByID(ctx context.Context, id string) (*ServerInvite, error)
	GetInviteByToken(ctx context.Context, token string) (*ServerInvite, error)
	ListInvites(ctx context.Context, limit, offset int) ([]*ServerInvite, error)
	UpdateInvite(ctx context.Context, invite *ServerInvite) error
}
