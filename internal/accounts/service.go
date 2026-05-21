package accounts

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

type FollowListener interface {
	OnFollow(ctx context.Context, followerID, followingID string, status string)
}

type Service struct {
	repo      Repository
	listeners []FollowListener
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:      repo,
		listeners: make([]FollowListener, 0),
	}
}

func (s *Service) RegisterFollowListener(l FollowListener) {
	s.listeners = append(s.listeners, l)
}

func (s *Service) GetByID(ctx context.Context, id string) (*Account, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByUsername(ctx context.Context, username string) (*Account, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *Service) UpdateProfile(ctx context.Context, accountID string, params *UpdateProfileParams) (*Account, error) {
	return s.repo.UpdateProfile(ctx, accountID, params)
}

func (s *Service) SearchByUsername(ctx context.Context, query string, limit int) ([]*Account, error) {
	if limit <= 0 || limit > 40 {
		limit = 40
	}
	return s.repo.SearchByUsername(ctx, query, limit)
}

func (s *Service) Follow(ctx context.Context, followerID, followingID string) (*Relationship, error) {
	if followerID == followingID {
		return nil, fmt.Errorf("cannot follow yourself")
	}

	blocked, _ := s.repo.IsBlocking(ctx, followingID, followerID)
	if blocked {
		return nil, fmt.Errorf("you are blocked by this user")
	}
	blocked, _ = s.repo.IsBlocking(ctx, followerID, followingID)
	if blocked {
		return nil, fmt.Errorf("unblock this user before following")
	}

	already, _ := s.repo.IsFollowing(ctx, followerID, followingID)
	if already {
		return s.GetRelationship(ctx, followerID, followingID)
	}

	target, err := s.repo.GetByID(ctx, followingID)
	if err != nil {
		return nil, fmt.Errorf("target account not found")
	}

	status := "accepted"
	if target.IsLocked {
		status = "pending"
	}

	if err := s.repo.Follow(ctx, followerID, followingID, status); err != nil {
		return nil, fmt.Errorf("creating follow: %w", err)
	}

	for _, l := range s.listeners {
		go l.OnFollow(context.Background(), followerID, followingID, status)
	}

	if status == "accepted" {
		if err := s.repo.IncrementFollowingCount(ctx, followerID); err != nil {
			log.Error().Err(err).Msg("failed to increment following count")
		}
		if err := s.repo.IncrementFollowersCount(ctx, followingID); err != nil {
			log.Error().Err(err).Msg("failed to increment followers count")
		}
	}

	return s.GetRelationship(ctx, followerID, followingID)
}

func (s *Service) Unfollow(ctx context.Context, followerID, followingID string) (*Relationship, error) {
	if followerID == followingID {
		return nil, fmt.Errorf("cannot unfollow yourself")
	}

	wasFollowing, _ := s.repo.IsFollowing(ctx, followerID, followingID)

	if err := s.repo.Unfollow(ctx, followerID, followingID); err != nil {
		return nil, fmt.Errorf("removing follow: %w", err)
	}

	if wasFollowing {
		if err := s.repo.DecrementFollowingCount(ctx, followerID); err != nil {
			log.Error().Err(err).Msg("failed to decrement following count")
		}
		if err := s.repo.DecrementFollowersCount(ctx, followingID); err != nil {
			log.Error().Err(err).Msg("failed to decrement followers count")
		}
	}

	return s.GetRelationship(ctx, followerID, followingID)
}

func (s *Service) AcceptFollowRequest(ctx context.Context, accountID, followerID string) error {
	if err := s.repo.AcceptFollow(ctx, followerID, accountID); err != nil {
		return fmt.Errorf("accepting follow: %w", err)
	}

	if err := s.repo.IncrementFollowingCount(ctx, followerID); err != nil {
		log.Error().Err(err).Msg("failed to increment following count")
	}
	if err := s.repo.IncrementFollowersCount(ctx, accountID); err != nil {
		log.Error().Err(err).Msg("failed to increment followers count")
	}
	return nil
}

func (s *Service) RejectFollowRequest(ctx context.Context, accountID, followerID string) error {
	return s.repo.RejectFollow(ctx, followerID, accountID)
}

func (s *Service) Block(ctx context.Context, accountID, targetID string) (*Relationship, error) {
	if accountID == targetID {
		return nil, fmt.Errorf("cannot block yourself")
	}

	wasFollowing, _ := s.repo.IsFollowing(ctx, accountID, targetID)
	wasFollowedBy, _ := s.repo.IsFollowing(ctx, targetID, accountID)

	s.repo.Unfollow(ctx, accountID, targetID)
	s.repo.Unfollow(ctx, targetID, accountID)

	if wasFollowing {
		s.repo.DecrementFollowingCount(ctx, accountID)
		s.repo.DecrementFollowersCount(ctx, targetID)
	}
	if wasFollowedBy {
		s.repo.DecrementFollowingCount(ctx, targetID)
		s.repo.DecrementFollowersCount(ctx, accountID)
	}

	if err := s.repo.Block(ctx, accountID, targetID); err != nil {
		return nil, fmt.Errorf("blocking: %w", err)
	}

	return s.GetRelationship(ctx, accountID, targetID)
}

func (s *Service) Unblock(ctx context.Context, accountID, targetID string) (*Relationship, error) {
	if err := s.repo.Unblock(ctx, accountID, targetID); err != nil {
		return nil, fmt.Errorf("unblocking: %w", err)
	}
	return s.GetRelationship(ctx, accountID, targetID)
}

func (s *Service) Mute(ctx context.Context, accountID, targetID string, hideNotifications bool) (*Relationship, error) {
	if accountID == targetID {
		return nil, fmt.Errorf("cannot mute yourself")
	}
	if err := s.repo.Mute(ctx, accountID, targetID, hideNotifications); err != nil {
		return nil, fmt.Errorf("muting: %w", err)
	}
	return s.GetRelationship(ctx, accountID, targetID)
}

func (s *Service) Unmute(ctx context.Context, accountID, targetID string) (*Relationship, error) {
	if err := s.repo.Unmute(ctx, accountID, targetID); err != nil {
		return nil, fmt.Errorf("unmuting: %w", err)
	}
	return s.GetRelationship(ctx, accountID, targetID)
}

func (s *Service) GetRelationship(ctx context.Context, accountID, targetID string) (*Relationship, error) {
	following, _ := s.repo.IsFollowing(ctx, accountID, targetID)
	followedBy, _ := s.repo.IsFollowing(ctx, targetID, accountID)
	blocking, _ := s.repo.IsBlocking(ctx, accountID, targetID)
	muting, _ := s.repo.IsMuting(ctx, accountID, targetID)

	return &Relationship{
		ID:         targetID,
		Following:  following,
		FollowedBy: followedBy,
		Blocking:   blocking,
		Muting:     muting,
	}, nil
}

func (s *Service) ListFollowers(ctx context.Context, accountID string, limit, offset int) ([]*Account, error) {
	return s.repo.ListFollowers(ctx, accountID, limit, offset)
}

func (s *Service) ListFollowing(ctx context.Context, accountID string, limit, offset int) ([]*Account, error) {
	return s.repo.ListFollowing(ctx, accountID, limit, offset)
}

func (s *Service) SetSuspended(ctx context.Context, id string, suspended bool) error {
	return s.repo.SetSuspended(ctx, id, suspended)
}

func (s *Service) SetSilenced(ctx context.Context, id string, silenced bool) error {
	return s.repo.SetSilenced(ctx, id, silenced)
}

func (s *Service) SetRole(ctx context.Context, id string, role string) error {
	if role != "admin" && role != "moderator" && role != "user" {
		return fmt.Errorf("invalid role")
	}
	return s.repo.SetRole(ctx, id, role)
}

func (s *Service) ListLocal(ctx context.Context, limit, offset int) ([]*Account, error) {
	return s.repo.ListLocal(ctx, limit, offset)
}
