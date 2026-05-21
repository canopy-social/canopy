package posts

import (
	"context"
	"time"
)

type Post struct {
	ID             string    `json:"id"`
	URI            string    `json:"uri"`
	URL            *string   `json:"url,omitempty"`
	AccountID      string    `json:"account_id"`
	Content        string    `json:"content"`
	ContentText    string    `json:"-"`
	ContentWarning *string   `json:"content_warning,omitempty"`
	IsSensitive    bool      `json:"is_sensitive"`
	Visibility     string    `json:"visibility"`
	Language       *string   `json:"language,omitempty"`
	ReplyToID      *string   `json:"reply_to_id,omitempty"`
	ThreadRootID   *string   `json:"thread_root_id,omitempty"`
	BoostOfID      *string   `json:"boost_of_id,omitempty"`
	IsLocal        bool      `json:"is_local"`
	IsPinned       bool      `json:"is_pinned"`
	LikesCount     int       `json:"likes_count"`
	BoostsCount    int       `json:"boosts_count"`
	RepliesCount   int       `json:"replies_count"`
	PostStyleID    *string   `json:"post_style_id,omitempty"`
	Tags           []string  `json:"tags,omitempty"`
	Mentions       []string  `json:"mentions,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreatePostParams struct {
	Content        string  `json:"content"`
	ContentWarning *string `json:"content_warning,omitempty"`
	IsSensitive    bool    `json:"is_sensitive"`
	Visibility     string  `json:"visibility"`
	Language       *string `json:"language,omitempty"`
	ReplyToID      *string `json:"reply_to_id,omitempty"`
}

type UpdatePostParams struct {
	Content        *string `json:"content,omitempty"`
	ContentWarning *string `json:"content_warning,omitempty"`
	IsSensitive    *bool   `json:"is_sensitive,omitempty"`
}

type ThreadContext struct {
	Ancestors   []*Post `json:"ancestors"`
	Descendants []*Post `json:"descendants"`
}

type Repository interface {
	Create(ctx context.Context, post *Post) (*Post, error)
	GetByID(ctx context.Context, id string) (*Post, error)
	GetByURI(ctx context.Context, uri string) (*Post, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, params *UpdatePostParams) (*Post, error)
	Pin(ctx context.Context, postID, accountID string) error
	Unpin(ctx context.Context, postID, accountID string) error

	ListByAccount(ctx context.Context, accountID string, includeBoosts bool, limit, offset int) ([]*Post, error)
	ListPinned(ctx context.Context, accountID string) ([]*Post, error)
	ListPublicTimeline(ctx context.Context, limit, offset int) ([]*Post, error)
	ListReplies(ctx context.Context, postID string) ([]*Post, error)
	ListThreadContext(ctx context.Context, threadRootID string) ([]*Post, error)

	Like(ctx context.Context, postID, accountID string) error
	Unlike(ctx context.Context, postID, accountID string) error
	HasLiked(ctx context.Context, postID, accountID string) (bool, error)
	Boost(ctx context.Context, postID, accountID string) error
	Unboost(ctx context.Context, postID, accountID string) error
	HasBoosted(ctx context.Context, postID, accountID string) (bool, error)

	AddMention(ctx context.Context, postID, accountID, uri string) error
	AddTag(ctx context.Context, postID, tag string) error
	SearchByTag(ctx context.Context, tag string, limit, offset int) ([]*Post, error)
}
