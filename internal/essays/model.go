package essays

import (
	"context"
	"time"
)

// Essay is the domain model for a long-form post (AP Article).
type Essay struct {
	ID                 string     `json:"id"`
	URI                string     `json:"uri"`
	URL                *string    `json:"url,omitempty"`
	AccountID          string     `json:"account_id"`
	Title              string     `json:"title"`
	Slug               string     `json:"slug"`
	Subtitle           *string    `json:"subtitle,omitempty"`
	Content            string     `json:"content"`
	ContentText        string     `json:"-"`
	ContentRaw         string     `json:"-"`
	CoverMediaID       *string    `json:"cover_media_id,omitempty"`
	ReadingTimeMinutes *int       `json:"reading_time_minutes,omitempty"`
	Visibility         string     `json:"visibility"`
	Language           *string    `json:"language,omitempty"`
	IsLocal            bool       `json:"is_local"`
	WordCount          int        `json:"word_count"`
	LikesCount         int        `json:"likes_count"`
	ViewsCount         int        `json:"views_count"`
	PublishedAt        *time.Time `json:"published_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// CreateEssayParams holds input for creating a new essay.
type CreateEssayParams struct {
	Title      string  `json:"title"`
	Subtitle   *string `json:"subtitle,omitempty"`
	Content    string  `json:"content"`
	ContentRaw string  `json:"content_raw"`
	Visibility string  `json:"visibility"`
	Language   *string `json:"language,omitempty"`
	Publish    bool    `json:"publish"`
}

// UpdateEssayParams holds input for updating an essay.
type UpdateEssayParams struct {
	Title      *string `json:"title,omitempty"`
	Subtitle   *string `json:"subtitle,omitempty"`
	Content    *string `json:"content,omitempty"`
	ContentRaw *string `json:"content_raw,omitempty"`
	Visibility *string `json:"visibility,omitempty"`
	Language   *string `json:"language,omitempty"`
}

// Repository defines data access for essays.
type Repository interface {
	Create(ctx context.Context, essay *Essay) (*Essay, error)
	GetByID(ctx context.Context, id string) (*Essay, error)
	GetBySlug(ctx context.Context, accountID, slug string) (*Essay, error)
	Update(ctx context.Context, id string, params *UpdateEssayParams) (*Essay, error)
	Publish(ctx context.Context, id string) (*Essay, error)
	Unpublish(ctx context.Context, id string) (*Essay, error)
	Delete(ctx context.Context, id string) error
	ListByAccount(ctx context.Context, accountID string, limit, offset int) ([]*Essay, error)
	ListDrafts(ctx context.Context, accountID string, limit, offset int) ([]*Essay, error)
	IncrementViews(ctx context.Context, id string) error
}
