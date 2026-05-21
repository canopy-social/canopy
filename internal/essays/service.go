package essays

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/microcosm-cc/bluemonday"

	"github.com/sumi-devs/canopy-social/canopy/pkg/config"
	"github.com/sumi-devs/canopy-social/canopy/pkg/ulid"
)

var (
	sanitizer   = bluemonday.UGCPolicy()
	slugCleanRe = regexp.MustCompile(`[^a-z0-9]+`)
)

// Service handles business logic for essays.
type Service struct {
	repo Repository
	cfg  *config.Config
}

// NewService creates a new essay service.
func NewService(repo Repository, cfg *config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

// Create creates a new essay.
func (s *Service) Create(ctx context.Context, accountID string, params *CreateEssayParams) (*Essay, error) {
	if strings.TrimSpace(params.Title) == "" {
		return nil, fmt.Errorf("title is required")
	}
	if strings.TrimSpace(params.Content) == "" {
		return nil, fmt.Errorf("content is required")
	}

	sanitizedContent := sanitizer.Sanitize(params.Content)
	plaintext := stripHTML(sanitizedContent)
	wordCount := countWords(plaintext)
	readingTime := (wordCount + 249) / 250 // ~250 wpm

	slug := generateSlug(params.Title)
	essayID := ulid.New()
	baseURL := s.cfg.BaseURL()
	essayURI := fmt.Sprintf("%s/essays/%s/%s", baseURL, accountID, slug)
	essayURL := essayURI

	essay := &Essay{
		ID:                 essayID,
		URI:                essayURI,
		URL:                &essayURL,
		AccountID:          accountID,
		Title:              params.Title,
		Slug:               slug,
		Subtitle:           params.Subtitle,
		Content:            sanitizedContent,
		ContentText:        plaintext,
		ContentRaw:         params.ContentRaw,
		Visibility:         params.Visibility,
		Language:           params.Language,
		IsLocal:            true,
		WordCount:          wordCount,
		ReadingTimeMinutes: &readingTime,
	}

	created, err := s.repo.Create(ctx, essay)
	if err != nil {
		return nil, fmt.Errorf("creating essay: %w", err)
	}

	if params.Publish {
		created, err = s.repo.Publish(ctx, created.ID)
		if err != nil {
			return nil, fmt.Errorf("publishing essay: %w", err)
		}
	}

	return created, nil
}

// GetByID returns an essay by ID, incrementing view count.
func (s *Service) GetByID(ctx context.Context, id string) (*Essay, error) {
	essay, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// Fire-and-forget view increment
	go s.repo.IncrementViews(context.Background(), id)
	return essay, nil
}

// GetBySlug returns an essay by account ID and slug.
func (s *Service) GetBySlug(ctx context.Context, accountID, slug string) (*Essay, error) {
	essay, err := s.repo.GetBySlug(ctx, accountID, slug)
	if err != nil {
		return nil, err
	}
	go s.repo.IncrementViews(context.Background(), essay.ID)
	return essay, nil
}

// Update updates an essay.
func (s *Service) Update(ctx context.Context, essayID, accountID string, params *UpdateEssayParams) (*Essay, error) {
	essay, err := s.repo.GetByID(ctx, essayID)
	if err != nil {
		return nil, fmt.Errorf("essay not found")
	}
	if essay.AccountID != accountID {
		return nil, fmt.Errorf("not authorized")
	}
	return s.repo.Update(ctx, essayID, params)
}

// Publish publishes a draft essay.
func (s *Service) Publish(ctx context.Context, essayID, accountID string) (*Essay, error) {
	essay, err := s.repo.GetByID(ctx, essayID)
	if err != nil {
		return nil, fmt.Errorf("essay not found")
	}
	if essay.AccountID != accountID {
		return nil, fmt.Errorf("not authorized")
	}
	return s.repo.Publish(ctx, essayID)
}

// Unpublish unpublishes a published essay (returns to draft).
func (s *Service) Unpublish(ctx context.Context, essayID, accountID string) (*Essay, error) {
	essay, err := s.repo.GetByID(ctx, essayID)
	if err != nil {
		return nil, fmt.Errorf("essay not found")
	}
	if essay.AccountID != accountID {
		return nil, fmt.Errorf("not authorized")
	}
	return s.repo.Unpublish(ctx, essayID)
}

// Delete deletes an essay.
func (s *Service) Delete(ctx context.Context, essayID, accountID string) error {
	essay, err := s.repo.GetByID(ctx, essayID)
	if err != nil {
		return fmt.Errorf("essay not found")
	}
	if essay.AccountID != accountID {
		return fmt.Errorf("not authorized")
	}
	return s.repo.Delete(ctx, essayID)
}

// ListByAccount returns published essays for an account.
func (s *Service) ListByAccount(ctx context.Context, accountID string, limit, offset int) ([]*Essay, error) {
	return s.repo.ListByAccount(ctx, accountID, limit, offset)
}

// ListDrafts returns draft essays for an account.
func (s *Service) ListDrafts(ctx context.Context, accountID string, limit, offset int) ([]*Essay, error) {
	return s.repo.ListDrafts(ctx, accountID, limit, offset)
}

// --- Helpers ---

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = slugCleanRe.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if utf8.RuneCountInString(slug) > 80 {
		slug = slug[:80]
	}
	if slug == "" {
		slug = ulid.New()[:8]
	}
	return slug
}

func countWords(s string) int {
	return len(strings.Fields(s))
}

func stripHTML(s string) string {
	p := bluemonday.StrictPolicy()
	return strings.TrimSpace(p.Sanitize(s))
}
