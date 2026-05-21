package timeline

import (
	"context"
	"github.com/sumi-devs/canopy-social/canopy/internal/accounts"
	"github.com/sumi-devs/canopy-social/canopy/internal/posts"
)

type PostRepository interface {
	GetByID(ctx context.Context, id string) (*posts.Post, error)
	ListPostsByAccountWithBoosts(ctx context.Context, accountID string, limit, offset int) ([]*posts.Post, error)
	ListPublicTimeline(ctx context.Context, local bool, limit, offset int) ([]*posts.Post, error)
}

type AccountRepository interface {
	GetByID(ctx context.Context, id string) (*accounts.Account, error)
	ListFollowers(ctx context.Context, accountID string, limit, offset int) ([]*accounts.Account, error)
	ListFollowing(ctx context.Context, accountID string, limit, offset int) ([]*accounts.Account, error)
}

type TimelineResponse struct {
	Data       []*posts.Post `json:"data"`
	NextCursor string        `json:"next_cursor"`
	PrevCursor string        `json:"prev_cursor"`
}
