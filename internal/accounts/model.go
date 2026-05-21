package accounts

import (
	"context"
	"time"
)

type Account struct {
	ID              string     `json:"id"`
	Username        string     `json:"username"`
	Domain          *string    `json:"domain,omitempty"`
	URI             string     `json:"uri"`
	DisplayName     *string    `json:"display_name,omitempty"`
	Bio             *string    `json:"bio,omitempty"`
	BioText         *string    `json:"bio_text,omitempty"`
	AvatarURL       *string    `json:"avatar_url,omitempty"`
	HeaderURL       *string    `json:"header_url,omitempty"`
	Role            string     `json:"role"`
	IsLocal         bool       `json:"is_local"`
	IsLocked        bool       `json:"is_locked"`
	IsBot           bool       `json:"is_bot"`
	IsSuspended     bool       `json:"-"`
	IsSilenced      bool       `json:"-"`
	ActorType       string     `json:"actor_type"`
	FollowersCount  int        `json:"followers_count"`
	FollowingCount  int        `json:"following_count"`
	PostsCount      int        `json:"posts_count"`
	PublicKeyPEM    string     `json:"-"`
	PrivateKeyPEM   *string    `json:"-"`
	KeyID           string     `json:"-"`
	InboxURL        *string    `json:"-"`
	OutboxURL       *string    `json:"-"`
	SharedInboxURL  *string    `json:"-"`
	FollowersURL    *string    `json:"-"`
	FollowingURL    *string    `json:"-"`
	PasswordHash    *string    `json:"-"`
	Email           *string    `json:"-"`
	EmailVerifiedAt *time.Time `json:"-"`
	CustomDomain    *string    `json:"custom_domain,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type Relationship struct {
	ID         string `json:"id"`
	Following  bool   `json:"following"`
	FollowedBy bool   `json:"followed_by"`
	Blocking   bool   `json:"blocking"`
	Muting     bool   `json:"muting"`
	Requested  bool   `json:"requested"`
}

type UpdateProfileParams struct {
	DisplayName *string `json:"display_name,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	IsLocked    *bool   `json:"is_locked,omitempty"`
	IsBot       *bool   `json:"is_bot,omitempty"`
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*Account, error)
	GetByURI(ctx context.Context, uri string) (*Account, error)
	GetByUsername(ctx context.Context, username string) (*Account, error)
	GetByUsernameAndDomain(ctx context.Context, username, domain string) (*Account, error)
	GetByEmail(ctx context.Context, email string) (*Account, error)
	Create(ctx context.Context, account *Account) (*Account, error)
	UpdateProfile(ctx context.Context, id string, params *UpdateProfileParams) (*Account, error)
	ListLocal(ctx context.Context, limit, offset int) ([]*Account, error)
	SearchByUsername(ctx context.Context, query string, limit int) ([]*Account, error)

	Follow(ctx context.Context, followerID, followingID, status string) error
	Unfollow(ctx context.Context, followerID, followingID string) error
	AcceptFollow(ctx context.Context, followerID, followingID string) error
	RejectFollow(ctx context.Context, followerID, followingID string) error
	IsFollowing(ctx context.Context, followerID, followingID string) (bool, error)
	ListFollowers(ctx context.Context, accountID string, limit, offset int) ([]*Account, error)
	ListFollowing(ctx context.Context, accountID string, limit, offset int) ([]*Account, error)
	ListPendingRequests(ctx context.Context, accountID string, limit, offset int) ([]*Account, error)

	Block(ctx context.Context, accountID, targetID string) error
	Unblock(ctx context.Context, accountID, targetID string) error
	IsBlocking(ctx context.Context, accountID, targetID string) (bool, error)
	Mute(ctx context.Context, accountID, targetID string, hideNotifications bool) error
	Unmute(ctx context.Context, accountID, targetID string) error
	IsMuting(ctx context.Context, accountID, targetID string) (bool, error)
	ListBlocks(ctx context.Context, accountID string, limit, offset int) ([]*Account, error)
	ListMutes(ctx context.Context, accountID string, limit, offset int) ([]*Account, error)

	IncrementFollowersCount(ctx context.Context, id string) error
	DecrementFollowersCount(ctx context.Context, id string) error
	IncrementFollowingCount(ctx context.Context, id string) error
	DecrementFollowingCount(ctx context.Context, id string) error
}
