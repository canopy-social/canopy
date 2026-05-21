package themes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/sumi-devs/canopy-social/canopy/pkg/ulid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func DefaultColorPalette() ColorPalette {
	return ColorPalette{
		Background:    "#ffffff",
		Surface:       "#f8f8f8",
		TextPrimary:   "#111111",
		TextSecondary: "#555555",
		Accent:        "#0066ff",
		AccentText:    "#ffffff",
		Border:        "#e0e0e0",
		Link:          "#0066ff",
	}
}

func DefaultFontPalette() FontPalette {
	return FontPalette{
		Body:    FontConfig{Family: "Inter", Size: 16, Weight: 400, LineHeight: 1.6, LetterSpacing: 0},
		Heading: FontConfig{Family: "Outfit", Size: 24, Weight: 700, LineHeight: 1.2, LetterSpacing: -0.02},
		Mono:    FontConfig{Family: "JetBrains Mono", Size: 14, Weight: 400, LineHeight: 1.5, LetterSpacing: 0},
		Display: FontConfig{Family: "Outfit", Size: 48, Weight: 900, LineHeight: 1.0, LetterSpacing: -0.03},
	}
}

func (s *Service) CreatePageTheme(ctx context.Context, accountID string) (*PageTheme, error) {
	existing, _ := s.repo.GetPageThemeByAccountID(ctx, accountID)
	if existing != nil {
		return nil, errors.New("account already has a page theme")
	}

	theme := &PageTheme{
		ID:                 ulid.New(),
		AccountID:          accountID,
		Colors:             DefaultColorPalette(),
		Fonts:              DefaultFontPalette(),
		Layout:             []LayoutWidget{},
		Stickers:           []Sticker{},
		Widgets:            []LayoutWidget{},
		BGType:             "color",
		BGImageSize:        "cover",
		BGBlur:             0,
		BGOpacity:          100,
		PageMaxWidth:       800,
		PagePadding:        24,
		ShowFollowerCount:  true,
		ShowFollowingCount: true,
		GardenMode:         false,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	css := GeneratePageCSS(theme)
	theme.GeneratedCSS = &css
	now := time.Now()
	theme.CSSGeneratedAt = &now

	return s.repo.CreatePageTheme(ctx, theme)
}

func (s *Service) GetPageTheme(ctx context.Context, accountID string) (*PageTheme, error) {
	return s.repo.GetPageThemeByAccountID(ctx, accountID)
}

func (s *Service) UpdatePageTheme(ctx context.Context, accountID string, update *PageThemeUpdateRequest) (*PageTheme, error) {
	theme, err := s.repo.GetPageThemeByAccountID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("theme not found: %w", err)
	}

	if update.Colors != nil {
		theme.Colors = *update.Colors
	}
	if update.Fonts != nil {
		theme.Fonts = *update.Fonts
	}
	if update.BGType != nil {
		theme.BGType = *update.BGType
	}
	if update.BGGradient != nil {
		theme.BGGradient = update.BGGradient
	}
	if update.BGImageID != nil {
		theme.BGImageID = update.BGImageID
	}
	if update.BGImageSize != nil {
		theme.BGImageSize = *update.BGImageSize
	}
	if update.BGBlur != nil {
		theme.BGBlur = *update.BGBlur
	}
	if update.BGOpacity != nil {
		theme.BGOpacity = *update.BGOpacity
	}
	if update.PageMaxWidth != nil {
		theme.PageMaxWidth = *update.PageMaxWidth
	}
	if update.PagePadding != nil {
		theme.PagePadding = *update.PagePadding
	}
	if update.ShowFollowerCount != nil {
		theme.ShowFollowerCount = *update.ShowFollowerCount
	}
	if update.ShowFollowingCount != nil {
		theme.ShowFollowingCount = *update.ShowFollowingCount
	}
	if update.GardenMode != nil {
		theme.GardenMode = *update.GardenMode
	}
	if update.InheritsServerTheme != nil {
		theme.InheritsServerTheme = *update.InheritsServerTheme
	}
	if update.ParentThemeID != nil {
		theme.ParentThemeID = update.ParentThemeID
	}
	if update.Layout != nil {
		theme.Layout = update.Layout
	}
	if update.Stickers != nil {
		theme.Stickers = update.Stickers
	}
	if update.Widgets != nil {
		theme.Widgets = update.Widgets
	}

	if err := ValidatePageTheme(theme); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	css := GeneratePageCSS(theme)
	theme.GeneratedCSS = &css
	now := time.Now()
	theme.CSSGeneratedAt = &now
	theme.UpdatedAt = now

	if err := s.repo.UpdatePageTheme(ctx, theme); err != nil {
		return nil, err
	}

	if err := s.autoSaveVersion(ctx, theme); err != nil {
		return nil, fmt.Errorf("auto-save failed: %w", err)
	}

	return theme, nil
}

func (s *Service) autoSaveVersion(ctx context.Context, theme *PageTheme) error {
	count, err := s.repo.CountAutoSavedVersions(ctx, theme.AccountID)
	if err != nil {
		return err
	}

	if count >= 50 {
		if err := s.repo.DeleteOldestAutoSave(ctx, theme.AccountID); err != nil {
			return err
		}
	}

	snapshot, err := json.Marshal(theme)
	if err != nil {
		return err
	}

	version := &ThemeVersion{
		ID:            ulid.New(),
		AccountID:     theme.AccountID,
		ThemeSnapshot: snapshot,
		AutoSaved:     true,
		CreatedAt:     time.Now(),
	}

	_, err = s.repo.CreateThemeVersion(ctx, version)
	return err
}

func (s *Service) SaveNamedVersion(ctx context.Context, accountID, label string) (*ThemeVersion, error) {
	theme, err := s.repo.GetPageThemeByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	count, err := s.repo.CountNamedVersions(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if count >= 20 {
		return nil, errors.New("maximum 20 named versions reached")
	}

	snapshot, err := json.Marshal(theme)
	if err != nil {
		return nil, err
	}

	version := &ThemeVersion{
		ID:            ulid.New(),
		AccountID:     accountID,
		ThemeSnapshot: snapshot,
		Label:         &label,
		AutoSaved:     false,
		CreatedAt:     time.Now(),
	}

	return s.repo.CreateThemeVersion(ctx, version)
}

func (s *Service) ListVersions(ctx context.Context, accountID string, limit, offset int) ([]*ThemeVersion, error) {
	return s.repo.ListThemeVersions(ctx, accountID, limit, offset)
}

func (s *Service) RestoreVersion(ctx context.Context, accountID, versionID string) (*PageTheme, error) {
	version, err := s.repo.GetThemeVersionByID(ctx, versionID)
	if err != nil {
		return nil, err
	}
	if version.AccountID != accountID {
		return nil, errors.New("version does not belong to this account")
	}

	theme, err := s.repo.GetPageThemeByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	var restored PageTheme
	if err := json.Unmarshal(version.ThemeSnapshot, &restored); err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot: %w", err)
	}

	restored.ID = theme.ID
	restored.AccountID = accountID
	restored.UpdatedAt = time.Now()

	css := GeneratePageCSS(&restored)
	restored.GeneratedCSS = &css
	now := time.Now()
	restored.CSSGeneratedAt = &now

	if err := s.repo.UpdatePageTheme(ctx, &restored); err != nil {
		return nil, err
	}

	return &restored, nil
}

func (s *Service) DeleteVersion(ctx context.Context, accountID, versionID string) error {
	version, err := s.repo.GetThemeVersionByID(ctx, versionID)
	if err != nil {
		return err
	}
	if version.AccountID != accountID {
		return errors.New("version does not belong to this account")
	}
	return s.repo.DeleteThemeVersion(ctx, versionID)
}

func (s *Service) CreatePostStyle(ctx context.Context, accountID string, req *CreatePostStyleRequest) (*PostStyle, error) {
	style := &PostStyle{
		ID:                ulid.New(),
		AccountID:         accountID,
		BackgroundColor:   req.BackgroundColor,
		BackgroundImageID: req.BackgroundImageID,
		TextColor:         req.TextColor,
		FontFamily:        req.FontFamily,
		FontSize:          req.FontSize,
		FontWeight:        req.FontWeight,
		BorderRadius:      req.BorderRadius,
		BorderColor:       req.BorderColor,
		BorderWidth:       req.BorderWidth,
		Padding:           req.Padding,
		HasTexture:        req.HasTexture,
		TextureType:       req.TextureType,
		CreatedAt:         time.Now(),
	}

	if err := ValidatePostStyle(style); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	created, err := s.repo.CreatePostStyle(ctx, style)
	if err != nil {
		return nil, err
	}

	css := GeneratePostStyleCSS(created)
	created.GeneratedCSS = &css

	return created, nil
}

func (s *Service) GetPostStyle(ctx context.Context, id string) (*PostStyle, error) {
	return s.repo.GetPostStyleByID(ctx, id)
}

func (s *Service) ListPostStyles(ctx context.Context, accountID string, limit, offset int) ([]*PostStyle, error) {
	return s.repo.ListPostStylesByAccount(ctx, accountID, limit, offset)
}

func (s *Service) DeletePostStyle(ctx context.Context, accountID, styleID string) error {
	style, err := s.repo.GetPostStyleByID(ctx, styleID)
	if err != nil {
		return err
	}
	if style.AccountID != accountID {
		return errors.New("post style does not belong to this account")
	}
	return s.repo.DeletePostStyle(ctx, styleID)
}

func (s *Service) GetOrCreateEssayTheme(ctx context.Context, essayID string) (*EssayTheme, error) {
	existing, err := s.repo.GetEssayThemeByEssayID(ctx, essayID)
	if err == nil {
		return existing, nil
	}

	theme := &EssayTheme{
		ID:        ulid.New(),
		EssayID:   essayID,
		Colors:    json.RawMessage("{}"),
		Fonts:     json.RawMessage("{}"),
		Layout:    json.RawMessage("[]"),
		BGType:    "inherit",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.CreateEssayTheme(ctx, theme)
}

func (s *Service) UpdateEssayTheme(ctx context.Context, essayID string, colors, fonts, layout json.RawMessage, bgType string) (*EssayTheme, error) {
	theme, err := s.repo.GetEssayThemeByEssayID(ctx, essayID)
	if err != nil {
		return nil, err
	}

	if colors != nil {
		theme.Colors = colors
	}
	if fonts != nil {
		theme.Fonts = fonts
	}
	if layout != nil {
		theme.Layout = layout
	}
	if bgType != "" {
		if err := ValidateBGType(bgType); err != nil {
			return nil, err
		}
		theme.BGType = bgType
	}

	theme.UpdatedAt = time.Now()
	if err := s.repo.UpdateEssayTheme(ctx, theme); err != nil {
		return nil, err
	}

	return theme, nil
}

func (s *Service) GetServerTheme(ctx context.Context) (*ServerTheme, error) {
	return s.repo.GetServerTheme(ctx)
}

func (s *Service) UpdateServerTheme(ctx context.Context, adminID string, colors, fonts, layout json.RawMessage, bgType string) (*ServerTheme, error) {
	theme, err := s.repo.GetServerTheme(ctx)
	if err != nil {
		theme = &ServerTheme{
			ID:     "singleton",
			Colors: json.RawMessage("{}"),
			Fonts:  json.RawMessage("{}"),
			Layout: json.RawMessage("[]"),
			BGType: "color",
		}
	}

	if colors != nil {
		theme.Colors = colors
	}
	if fonts != nil {
		theme.Fonts = fonts
	}
	if layout != nil {
		theme.Layout = layout
	}
	if bgType != "" {
		if err := ValidateBGType(bgType); err != nil {
			return nil, err
		}
		theme.BGType = bgType
	}

	theme.UpdatedBy = &adminID
	theme.UpdatedAt = time.Now()

	css := GenerateServerThemeCSS(theme)
	theme.GeneratedCSS = &css
	now := time.Now()
	theme.CSSGeneratedAt = &now

	if err := s.repo.UpsertServerTheme(ctx, theme); err != nil {
		return nil, err
	}

	return theme, nil
}

func (s *Service) GetThemeCSS(ctx context.Context, accountID string) (string, error) {
	theme, err := s.repo.GetPageThemeByAccountID(ctx, accountID)
	if err != nil {
		return "", err
	}
	if theme.GeneratedCSS != nil {
		return *theme.GeneratedCSS, nil
	}
	return GeneratePageCSS(theme), nil
}

type PageThemeUpdateRequest struct {
	Colors             *ColorPalette  `json:"colors,omitempty"`
	Fonts              *FontPalette   `json:"fonts,omitempty"`
	BGType             *string        `json:"bg_type,omitempty"`
	BGGradient         *BGGradient    `json:"bg_gradient,omitempty"`
	BGImageID          *string        `json:"bg_image_id,omitempty"`
	BGImageSize        *string        `json:"bg_image_size,omitempty"`
	BGBlur             *int           `json:"bg_blur,omitempty"`
	BGOpacity          *int           `json:"bg_opacity,omitempty"`
	PageMaxWidth       *int           `json:"page_max_width,omitempty"`
	PagePadding        *int           `json:"page_padding,omitempty"`
	ShowFollowerCount  *bool          `json:"show_follower_count,omitempty"`
	ShowFollowingCount *bool          `json:"show_following_count,omitempty"`
	GardenMode         *bool          `json:"garden_mode,omitempty"`
	InheritsServerTheme *bool         `json:"inherits_server_theme,omitempty"`
	ParentThemeID      *string        `json:"parent_theme_id,omitempty"`
	Layout             []LayoutWidget `json:"layout,omitempty"`
	Stickers           []Sticker      `json:"stickers,omitempty"`
	Widgets            []LayoutWidget `json:"widgets,omitempty"`
}

type CreatePostStyleRequest struct {
	BackgroundColor   *string `json:"background_color,omitempty"`
	BackgroundImageID *string `json:"background_image_id,omitempty"`
	TextColor         *string `json:"text_color,omitempty"`
	FontFamily        *string `json:"font_family,omitempty"`
	FontSize          *int    `json:"font_size,omitempty"`
	FontWeight        *int    `json:"font_weight,omitempty"`
	BorderRadius      *int    `json:"border_radius,omitempty"`
	BorderColor       *string `json:"border_color,omitempty"`
	BorderWidth       *int    `json:"border_width,omitempty"`
	Padding           *int    `json:"padding,omitempty"`
	HasTexture        bool    `json:"has_texture"`
	TextureType       *string `json:"texture_type,omitempty"`
}
