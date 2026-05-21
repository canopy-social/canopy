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

type PostCreatedListener interface {
	OnPostCreated(ctx context.Context, post *Post)
}

type PostLikedListener interface {
	OnPostLiked(ctx context.Context, postID, accountID string)
}

type PostBoostedListener interface {
	OnPostBoosted(ctx context.Context, postID, accountID string)
}

type Service struct {
	repo           Repository
	cfg            *config.Config
	listeners      []PostCreatedListener
	likeListeners  []PostLikedListener
	boostListeners []PostBoostedListener
}

func NewService(repo Repository, cfg *config.Config) *Service {
	return &Service{
		repo:           repo,
		cfg:            cfg,
		listeners:      make([]PostCreatedListener, 0),
		likeListeners:  make([]PostLikedListener, 0),
		boostListeners: make([]PostBoostedListener, 0),
	}
}

func (s *Service) RegisterListener(l PostCreatedListener) {
	s.listeners = append(s.listeners, l)
}

func (s *Service) RegisterLikeListener(l PostLikedListener) {
	s.likeListeners = append(s.likeListeners, l)
}

func (s *Service) RegisterBoostListener(l PostBoostedListener) {
	s.boostListeners = append(s.boostListeners, l)
}

func (s *Service) Create(ctx context.Context, accountID string, params *CreatePostParams) (*Post, error) {

	plaintext := stripHTML(params.Content)
	if len(plaintext) == 0 {
		return nil, fmt.Errorf("post content cannot be empty")
	}
	if len(plaintext) > s.cfg.Features.MaxPostLength {
		return nil, fmt.Errorf("post exceeds max length of %d characters", s.cfg.Features.MaxPostLength)
	}

	if !validate.Visibility(params.Visibility) {
		params.Visibility = "public"
	}

	sanitizedContent := sanitizer.Sanitize(params.Content)

	postID := ulid.New()
	baseURL := s.cfg.BaseURL()
	postURI := fmt.Sprintf("%s/posts/%s", baseURL, postID)
	postURL := postURI

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
		Language:       params.Language,
		ReplyToID:      params.ReplyToID,
		ThreadRootID:   threadRootID,
		IsLocal:        true,
	}

	created, err := s.repo.Create(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("creating post: %w", err)
	}

	if params.ReplyToID != nil {

		go func() {
			if err := s.repo.Like(context.Background(), *params.ReplyToID, ""); err != nil {

			}
		}()
	}

	mentions := mentionRe.FindAllStringSubmatch(params.Content, -1)
	for _, m := range mentions {
		username := m[1]
		if err := s.repo.AddMention(ctx, postID, username, ""); err != nil {
			log.Warn().Err(err).Str("mention", username).Msg("failed to store mention")
		}
	}

	tags := hashtagRe.FindAllStringSubmatch(params.Content, -1)
	for _, t := range tags {
		tag := strings.ToLower(t[1])
		if err := s.repo.AddTag(ctx, postID, tag); err != nil {
			log.Warn().Err(err).Str("tag", tag).Msg("failed to store tag")
		}
	}

	for _, l := range s.listeners {
		go l.OnPostCreated(context.Background(), created)
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*Post, error) {
	return s.repo.GetByID(ctx, id)
}

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

func (s *Service) Like(ctx context.Context, postID, accountID string) error {
	already, _ := s.repo.HasLiked(ctx, postID, accountID)
	if already {
		return nil
	}
	if err := s.repo.Like(ctx, postID, accountID); err != nil {
		return err
	}
	for _, l := range s.likeListeners {
		go l.OnPostLiked(context.Background(), postID, accountID)
	}
	return nil
}

func (s *Service) Unlike(ctx context.Context, postID, accountID string) error {
	return s.repo.Unlike(ctx, postID, accountID)
}

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
	if err := s.repo.Boost(ctx, postID, accountID); err != nil {
		return err
	}
	for _, l := range s.boostListeners {
		go l.OnPostBoosted(context.Background(), postID, accountID)
	}
	return nil
}

func (s *Service) Unboost(ctx context.Context, postID, accountID string) error {
	return s.repo.Unboost(ctx, postID, accountID)
}

func (s *Service) GetThreadContext(ctx context.Context, postID string) (*ThreadContext, error) {
	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}

	ctx2 := context.Background()
	var ancestors []*Post
	var descendants []*Post

	current := post
	for current.ReplyToID != nil {
		parent, err := s.repo.GetByID(ctx2, *current.ReplyToID)
		if err != nil {
			break
		}
		ancestors = append([]*Post{parent}, ancestors...)
		current = parent
	}

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

func (s *Service) ListByAccount(ctx context.Context, accountID string, includeBoosts bool, limit, offset int) ([]*Post, error) {
	return s.repo.ListByAccount(ctx, accountID, includeBoosts, limit, offset)
}

func (s *Service) ListPublicTimeline(ctx context.Context, limit, offset int) ([]*Post, error) {
	return s.repo.ListPublicTimeline(ctx, limit, offset)
}

func (s *Service) SearchByTag(ctx context.Context, tag string, limit, offset int) ([]*Post, error) {
	return s.repo.SearchByTag(ctx, strings.ToLower(tag), limit, offset)
}

func stripHTML(s string) string {
	p := bluemonday.StrictPolicy()
	return strings.TrimSpace(p.Sanitize(s))
}
