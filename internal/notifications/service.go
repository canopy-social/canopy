package notifications

import (
	"context"
	"regexp"
	"time"

	"github.com/sumi-devs/canopy-social/canopy/internal/accounts"
	"github.com/sumi-devs/canopy-social/canopy/internal/posts"
	"github.com/sumi-devs/canopy-social/canopy/pkg/ulid"
)

var mentionRe = regexp.MustCompile(`@([a-zA-Z0-9_]+)(?:@([a-zA-Z0-9.\-]+))?`)

type AccountLoader interface {
	GetByID(ctx context.Context, id string) (*accounts.Account, error)
	GetByUsername(ctx context.Context, username string) (*accounts.Account, error)
}

type PostLoader interface {
	GetByID(ctx context.Context, id string) (*posts.Post, error)
}

type Service struct {
	repo          Repository
	accountLoader AccountLoader
	postLoader    PostLoader
}

func NewService(repo Repository, accountLoader AccountLoader, postLoader PostLoader) *Service {
	return &Service{
		repo:          repo,
		accountLoader: accountLoader,
		postLoader:    postLoader,
	}
}

func (s *Service) OnFollow(ctx context.Context, followerID, followingID string, status string) {
	if followerID == followingID {
		return
	}
	n := &Notification{
		ID:            ulid.New(),
		AccountID:     followingID,
		Type:          "follow",
		FromAccountID: &followerID,
		CreatedAt:     time.Now(),
	}
	_, _ = s.repo.Create(ctx, n)
}

func (s *Service) OnPostLiked(ctx context.Context, postID, accountID string) {
	post, err := s.postLoader.GetByID(ctx, postID)
	if err != nil || post == nil {
		return
	}
	if post.AccountID == accountID {
		return
	}
	n := &Notification{
		ID:            ulid.New(),
		AccountID:     post.AccountID,
		Type:          "favourite",
		FromAccountID: &accountID,
		PostID:        &postID,
		CreatedAt:     time.Now(),
	}
	_, _ = s.repo.Create(ctx, n)
}

func (s *Service) OnPostBoosted(ctx context.Context, postID, accountID string) {
	post, err := s.postLoader.GetByID(ctx, postID)
	if err != nil || post == nil {
		return
	}
	if post.AccountID == accountID {
		return
	}
	n := &Notification{
		ID:            ulid.New(),
		AccountID:     post.AccountID,
		Type:          "reblog",
		FromAccountID: &accountID,
		PostID:        &postID,
		CreatedAt:     time.Now(),
	}
	_, _ = s.repo.Create(ctx, n)
}

func (s *Service) OnPostCreated(ctx context.Context, post *posts.Post) {
	notified := make(map[string]bool)
	if post.ReplyToID != nil {
		parent, err := s.postLoader.GetByID(ctx, *post.ReplyToID)
		if err == nil && parent != nil {
			if parent.AccountID != post.AccountID {
				n := &Notification{
					ID:            ulid.New(),
					AccountID:     parent.AccountID,
					Type:          "mention",
					FromAccountID: &post.AccountID,
					PostID:        &post.ID,
					CreatedAt:     time.Now(),
				}
				_, err = s.repo.Create(ctx, n)
				if err == nil {
					notified[parent.AccountID] = true
				}
			}
		}
	}
	mentions := mentionRe.FindAllStringSubmatch(post.ContentText, -1)
	for _, m := range mentions {
		username := m[1]
		acc, err := s.accountLoader.GetByUsername(ctx, username)
		if err == nil && acc != nil {
			if acc.ID != post.AccountID && !notified[acc.ID] {
				n := &Notification{
					ID:            ulid.New(),
					AccountID:     acc.ID,
					Type:          "mention",
					FromAccountID: &post.AccountID,
					PostID:        &post.ID,
					CreatedAt:     time.Now(),
				}
				_, err = s.repo.Create(ctx, n)
				if err == nil {
					notified[acc.ID] = true
				}
			}
		}
	}
}

func (s *Service) GetByID(ctx context.Context, id string) (*Notification, error) {
	n, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.enrichOne(ctx, n)
	return n, nil
}

func (s *Service) List(ctx context.Context, accountID string, limit int, maxID, sinceID string, types []string, excludeTypes []string) ([]*Notification, error) {
	list, err := s.repo.List(ctx, accountID, limit, maxID, sinceID, types, excludeTypes)
	if err != nil {
		return nil, err
	}
	for _, n := range list {
		s.enrichOne(ctx, n)
	}
	return list, nil
}

func (s *Service) MarkRead(ctx context.Context, accountID string, id string) (*Notification, error) {
	n, err := s.repo.MarkRead(ctx, accountID, id)
	if err != nil {
		return nil, err
	}
	s.enrichOne(ctx, n)
	return n, nil
}

func (s *Service) MarkAllRead(ctx context.Context, accountID string) error {
	return s.repo.MarkAllRead(ctx, accountID)
}

func (s *Service) Dismiss(ctx context.Context, accountID string, id string) error {
	return s.repo.Dismiss(ctx, accountID, id)
}

func (s *Service) enrichOne(ctx context.Context, n *Notification) {
	if n == nil {
		return
	}
	if n.FromAccountID != nil {
		acc, err := s.accountLoader.GetByID(ctx, *n.FromAccountID)
		if err == nil {
			n.Account = acc
		}
	}
	if n.PostID != nil {
		p, err := s.postLoader.GetByID(ctx, *n.PostID)
		if err == nil {
			n.Status = p
		}
	}
}
