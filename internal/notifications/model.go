package notifications

import (
	"context"
	"time"

	"github.com/sumi-devs/canopy-social/canopy/internal/accounts"
	"github.com/sumi-devs/canopy-social/canopy/internal/posts"
)

type Notification struct {
	ID            string            `json:"id"`
	AccountID     string            `json:"account_id"`
	Type          string            `json:"type"`
	FromAccountID *string           `json:"from_account_id,omitempty"`
	PostID        *string           `json:"post_id,omitempty"`
	EssayID       *string           `json:"essay_id,omitempty"`
	ChannelID     *string           `json:"channel_id,omitempty"`
	ReadAt        *time.Time        `json:"read_at,omitempty"`
	DismissedAt   *time.Time        `json:"dismissed_at,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	Account       *accounts.Account `json:"account,omitempty"`
	Status        *posts.Post       `json:"status,omitempty"`
}

type Repository interface {
	Create(ctx context.Context, n *Notification) (*Notification, error)
	GetByID(ctx context.Context, id string) (*Notification, error)
	List(ctx context.Context, accountID string, limit int, maxID, sinceID string, types []string, excludeTypes []string) ([]*Notification, error)
	MarkRead(ctx context.Context, accountID string, id string) (*Notification, error)
	MarkAllRead(ctx context.Context, accountID string) error
	Dismiss(ctx context.Context, accountID string, id string) error
}
