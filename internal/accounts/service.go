package accounts

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

// Service handles business logic for accounts.
type Service struct {
	repo Repository
}

// NewService creates a new account service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetByID returns an account by its ULID.
func (s *Service) GetByID(ctx context.Context, id string) (*Account, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByUsername returns a local account by username.
func (s *Service) GetByUsername(ctx context.Context, username string) (*Account, error) {
	return s.repo.GetByUsername(ctx, username)
}

// UpdateProfile updates a user's profile fields.
func (s *Service) UpdateProfile(ctx context.Context, accountID string, params *UpdateProfileParams) (*Account, error) {
	return s.repo.UpdateProfile(ctx, accountID, params)
}

// SearchByUsername searches for accounts by username prefix.
func (s *Service) SearchByUsername(ctx context.Context, query string, limit int) ([]*Account, error) {
	if limit <= 0 || limit > 40 {
		limit = 40
	}
	return s.repo.SearchByUsername(ctx, query, limit)
}

// Follow initiates a follow from followerID to followingID.
// If the target account is locked, the follow is set to "pending".
func (s *Service) Follow(ctx context.Context, followerID, followingID string) (*Relationship, error) {
	if followerID == followingID {
		return nil, fmt.Errorf("cannot follow yourself")
	}

	// Check block in both directions
	blocked, _ := s.repo.IsBlocking(ctx, followingID, followerID)
	if blocked {
		return nil, fmt.Errorf("you are blocked by this user")
	}
	blocked, _ = s.repo.IsBlocking(ctx, followerID, followingID)
	if blocked {
		return nil, fmt.Errorf("unblock this user before following")
	}

	// Check if already following
	already, _ := s.repo.IsFollowing(ctx, followerID, followingID)
	if already {
		return s.GetRelationship(ctx, followerID, followingID)
	}

	// Check if target is locked
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

	// Update counts only if immediately accepted
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

// Unfollow removes a follow relationship.
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

// AcceptFollowRequest accepts a pending follow request.
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

// RejectFollowRequest rejects a pending follow request.
func (s *Service) RejectFollowRequest(ctx context.Context, accountID, followerID string) error {
	return s.repo.RejectFollow(ctx, followerID, accountID)
}

// Block blocks a target account. Also removes any existing follow in both directions.
func (s *Service) Block(ctx context.Context, accountID, targetID string) (*Relationship, error) {
	if accountID == targetID {
		return nil, fmt.Errorf("cannot block yourself")
	}

	// Remove follow in both directions
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

// Unblock removes a block.
func (s *Service) Unblock(ctx context.Context, accountID, targetID string) (*Relationship, error) {
	if err := s.repo.Unblock(ctx, accountID, targetID); err != nil {
		return nil, fmt.Errorf("unblocking: %w", err)
	}
	return s.GetRelationship(ctx, accountID, targetID)
}

// Mute mutes a target account.
func (s *Service) Mute(ctx context.Context, accountID, targetID string, hideNotifications bool) (*Relationship, error) {
	if accountID == targetID {
		return nil, fmt.Errorf("cannot mute yourself")
	}
	if err := s.repo.Mute(ctx, accountID, targetID, hideNotifications); err != nil {
		return nil, fmt.Errorf("muting: %w", err)
	}
	return s.GetRelationship(ctx, accountID, targetID)
}

// Unmute removes a mute.
func (s *Service) Unmute(ctx context.Context, accountID, targetID string) (*Relationship, error) {
	if err := s.repo.Unmute(ctx, accountID, targetID); err != nil {
		return nil, fmt.Errorf("unmuting: %w", err)
	}
	return s.GetRelationship(ctx, accountID, targetID)
}

// GetRelationship returns the relationship between two accounts.
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

// ListFollowers returns paginated followers of an account.
func (s *Service) ListFollowers(ctx context.Context, accountID string, limit, offset int) ([]*Account, error) {
	return s.repo.ListFollowers(ctx, accountID, limit, offset)
}

// ListFollowing returns paginated accounts that an account follows.
func (s *Service) ListFollowing(ctx context.Context, accountID string, limit, offset int) ([]*Account, error) {
	return s.repo.ListFollowing(ctx, accountID, limit, offset)
}
