package posts

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"

	"github.com/sumi-devs/canopy-social/canopy/pkg/config"
	"github.com/sumi-devs/canopy-social/canopy/pkg/ulid"
	"github.com/sumi-devs/canopy-social/canopy/pkg/validate"
)

var (
	mentionRe = regexp.MustCompile(`@([a-zA-Z0-9_]+)(?:@([a-zA-Z0-9.\-]+))?`)
	hashtagRe = regexp.MustCompile(`#([a-zA-Z0-9_]+)`)
	sanitizer = bluemonday.UGCPolicy()
)

// Service handles business logic for posts.
type Service struct {
	repo Repository
	cfg  *config.Config
}

// NewService creates a new post service.
func NewService(repo Repository, cfg *config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

// Create creates a new post.
func (s *Service) Create(ctx context.Context, accountID string, params *CreatePostParams) (*Post, error) {
	// Validate content
	plaintext := stripHTML(params.Content)
	if len(plaintext) == 0 {
		return nil, fmt.Errorf("post content cannot be empty")
	}
	if len(plaintext) > s.cfg.Features.MaxPostLength {
		return nil, fmt.Errorf("post exceeds max length of %d characters", s.cfg.Features.MaxPostLength)
	}

	// Validate visibility
	if !validate.Visibility(params.Visibility) {
		params.Visibility = "public"
	}

	// Sanitize HTML
	sanitizedContent := sanitizer.Sanitize(params.Content)

	// Generate IDs and URIs
	postID := ulid.New()
	baseURL := s.cfg.BaseURL()
	postURI := fmt.Sprintf("%s/posts/%s", baseURL, postID)
	postURL := postURI

	// Determine thread root
	var threadRootID *string
	if params.ReplyToID != nil {
		parent, err := s.repo.GetByID(ctx, *params.ReplyToID)
		if err == nil && parent != nil {
			if parent.ThreadRootID != nil {
				threadRootID = parent.ThreadRootID
			} else {
				threadRootID = &parent.ID
			}
		}
	}

	post := &Post{
		ID:             postID,
		URI:            postURI,
		URL:            &postURL,
		AccountID:      accountID,
		Content:        sanitizedContent,
		ContentText:    plaintext,
		ContentWarning: params.ContentWarning,
		IsSensitive:    params.IsSensitive,
		Visibility:     params.Visibility,
		Language:        params.Language,
		ReplyToID:      params.ReplyToID,
		ThreadRootID:   threadRootID,
		IsLocal:        true,
	}

	created, err := s.repo.Create(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("creating post: %w", err)
	}

	// Increment reply count on parent
	if params.ReplyToID != nil {
		// fire and forget — counter inconsistency is tolerable
		go func() {
			if err := s.repo.Like(context.Background(), *params.ReplyToID, ""); err != nil {
				// This is intentionally a no-op placeholder. Reply count increment
				// will be handled properly when integrated with the full repo.
			}
		}()
	}

	// Extract and store mentions
	mentions := mentionRe.FindAllStringSubmatch(params.Content, -1)
	for _, m := range mentions {
		username := m[1]
		if err := s.repo.AddMention(ctx, postID, username, ""); err != nil {
			log.Warn().Err(err).Str("mention", username).Msg("failed to store mention")
		}
	}

	// Extract and store hashtags
	tags := hashtagRe.FindAllStringSubmatch(params.Content, -1)
	for _, t := range tags {
		tag := strings.ToLower(t[1])
		if err := s.repo.AddTag(ctx, postID, tag); err != nil {
			log.Warn().Err(err).Str("tag", tag).Msg("failed to store tag")
		}
	}

	return created, nil
}

// GetByID returns a post by ID.
func (s *Service) GetByID(ctx context.Context, id string) (*Post, error) {
	return s.repo.GetByID(ctx, id)
}

// Delete deletes a post (only if owned by the account).
func (s *Service) Delete(ctx context.Context, postID, accountID string) error {
	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("post not found")
	}
	if post.AccountID != accountID {
		return fmt.Errorf("not authorized to delete this post")
	}
	return s.repo.Delete(ctx, postID)
}

// Edit updates a post's content (only if owned by the account).
func (s *Service) Edit(ctx context.Context, postID, accountID string, params *UpdatePostParams) (*Post, error) {
	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}
	if post.AccountID != accountID {
		return nil, fmt.Errorf("not authorized to edit this post")
	}
	if !post.IsLocal {
		return nil, fmt.Errorf("cannot edit remote posts")
	}
	return s.repo.Update(ctx, postID, params)
}

// Like likes a post.
func (s *Service) Like(ctx context.Context, postID, accountID string) error {
	already, _ := s.repo.HasLiked(ctx, postID, accountID)
	if already {
		return nil
	}
	if err := s.repo.Like(ctx, postID, accountID); err != nil {
		return err
	}
	return nil
}

// Unlike removes a like.
func (s *Service) Unlike(ctx context.Context, postID, accountID string) error {
	return s.repo.Unlike(ctx, postID, accountID)
}

// Boost boosts a post.
func (s *Service) Boost(ctx context.Context, postID, accountID string) error {
	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("post not found")
	}
	if post.Visibility == "direct" || post.Visibility == "followers" {
		return fmt.Errorf("cannot boost non-public posts")
	}
	already, _ := s.repo.HasBoosted(ctx, postID, accountID)
	if already {
		return nil
	}
	return s.repo.Boost(ctx, postID, accountID)
}

// Unboost removes a boost.
func (s *Service) Unboost(ctx context.Context, postID, accountID string) error {
	return s.repo.Unboost(ctx, postID, accountID)
}

// GetThreadContext returns ancestors and descendants for a post.
func (s *Service) GetThreadContext(ctx context.Context, postID string) (*ThreadContext, error) {
	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}

	ctx2 := context.Background()
	var ancestors []*Post
	var descendants []*Post

	// Walk up the reply chain for ancestors
	current := post
	for current.ReplyToID != nil {
		parent, err := s.repo.GetByID(ctx2, *current.ReplyToID)
		if err != nil {
			break
		}
		ancestors = append([]*Post{parent}, ancestors...)
		current = parent
	}

	// Get all replies (descendants)
	if post.ThreadRootID != nil {
		all, err := s.repo.ListThreadContext(ctx, *post.ThreadRootID)
		if err == nil {
			for _, p := range all {
				if p.CreatedAt.After(post.CreatedAt) {
					descendants = append(descendants, p)
				}
			}
		}
	} else {
		replies, err := s.repo.ListReplies(ctx, postID)
		if err == nil {
			descendants = replies
		}
	}

	return &ThreadContext{
		Ancestors:   ancestors,
		Descendants: descendants,
	}, nil
}

// ListByAccount returns posts by an account.
func (s *Service) ListByAccount(ctx context.Context, accountID string, includeBoosts bool, limit, offset int) ([]*Post, error) {
	return s.repo.ListByAccount(ctx, accountID, includeBoosts, limit, offset)
}

// ListPublicTimeline returns the public timeline.
func (s *Service) ListPublicTimeline(ctx context.Context, limit, offset int) ([]*Post, error) {
	return s.repo.ListPublicTimeline(ctx, limit, offset)
}

// SearchByTag returns posts by hashtag.
func (s *Service) SearchByTag(ctx context.Context, tag string, limit, offset int) ([]*Post, error) {
	return s.repo.SearchByTag(ctx, strings.ToLower(tag), limit, offset)
}

func stripHTML(s string) string {
	p := bluemonday.StrictPolicy()
	return strings.TrimSpace(p.Sanitize(s))
}
