package media

import (
	"context"
	"time"
)

type Attachment struct {
	ID              string     `json:"id"`
	AccountID       string     `json:"account_id"`
	PostID          *string    `json:"post_id,omitempty"`
	EssayID         *string    `json:"essay_id,omitempty"`
	RemoteURL       *string    `json:"remote_url,omitempty"`
	URL             *string    `json:"url,omitempty"`
	ThumbnailURL    *string    `json:"thumbnail_url,omitempty"`
	Type            string     `json:"type"`
	MimeType        *string    `json:"mime_type,omitempty"`
	FileSize        *int64     `json:"file_size,omitempty"`
	Width           *int       `json:"width,omitempty"`
	Height          *int       `json:"height,omitempty"`
	DurationSeconds *float64   `json:"duration_seconds,omitempty"`
	Blurhash        *string    `json:"blurhash,omitempty"`
	AltText         *string    `json:"alt_text,omitempty"`
	IsProcessed     bool       `json:"is_processed"`
	ProcessingError *string    `json:"processing_error,omitempty"`
	StorageKey      *string    `json:"storage_key,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type Repository interface {
	Create(ctx context.Context, att *Attachment) (*Attachment, error)
	GetByID(ctx context.Context, id string) (*Attachment, error)
	Update(ctx context.Context, id string, altText *string) (*Attachment, error)
	Delete(ctx context.Context, id string) error
}
