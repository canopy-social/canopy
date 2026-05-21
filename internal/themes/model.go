package themes

import (
	"context"
	"encoding/json"
	"time"
)

type ColorPalette struct {
	Background    string `json:"background"`
	Surface       string `json:"surface"`
	TextPrimary   string `json:"text_primary"`
	TextSecondary string `json:"text_secondary"`
	Accent        string `json:"accent"`
	AccentText    string `json:"accent_text"`
	Border        string `json:"border"`
	Link          string `json:"link"`
}

type FontConfig struct {
	Family        string  `json:"family"`
	Size          int     `json:"size"`
	Weight        int     `json:"weight"`
	LineHeight    float64 `json:"line_height"`
	LetterSpacing float64 `json:"letter_spacing"`
}

type FontPalette struct {
	Body    FontConfig `json:"body"`
	Heading FontConfig `json:"heading"`
	Mono    FontConfig `json:"mono"`
	Display FontConfig `json:"display"`
}

type LayoutWidget struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	X        int             `json:"x"`
	Y        int             `json:"y"`
	Width    int             `json:"width"`
	Height   int             `json:"height"`
	Settings json.RawMessage `json:"settings,omitempty"`
}

type Sticker struct {
	ID      string `json:"id"`
	MediaID string `json:"media_id"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Rotate  int    `json:"rotate"`
	ZIndex  int    `json:"z_index"`
}

type GradientStop struct {
	Color    string  `json:"color"`
	Position float64 `json:"position"`
}

type BGGradient struct {
	Type   string         `json:"type"`
	Angle  int            `json:"angle"`
	Stops  []GradientStop `json:"stops"`
}

type PageTheme struct {
	ID                  string          `json:"id"`
	AccountID           string          `json:"account_id"`
	Colors              ColorPalette    `json:"colors"`
	Fonts               FontPalette     `json:"fonts"`
	Layout              []LayoutWidget  `json:"layout"`
	Stickers            []Sticker       `json:"stickers"`
	Widgets             []LayoutWidget  `json:"widgets"`
	BGType              string          `json:"bg_type"`
	BGGradient          *BGGradient     `json:"bg_gradient,omitempty"`
	BGImageID           *string         `json:"bg_image_id,omitempty"`
	BGImageSize         string          `json:"bg_image_size"`
	BGBlur              int             `json:"bg_blur"`
	BGOpacity           int             `json:"bg_opacity"`
	PageMaxWidth        int             `json:"page_max_width"`
	PagePadding         int             `json:"page_padding"`
	ShowFollowerCount   bool            `json:"show_follower_count"`
	ShowFollowingCount  bool            `json:"show_following_count"`
	GardenMode          bool            `json:"garden_mode"`
	InheritsServerTheme bool            `json:"inherits_server_theme"`
	ParentThemeID       *string         `json:"parent_theme_id,omitempty"`
	GeneratedCSS        *string         `json:"generated_css,omitempty"`
	CSSGeneratedAt      *time.Time      `json:"css_generated_at,omitempty"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

type ThemeVersion struct {
	ID            string          `json:"id"`
	AccountID     string          `json:"account_id"`
	ThemeSnapshot json.RawMessage `json:"theme_snapshot"`
	Label         *string         `json:"label,omitempty"`
	AutoSaved     bool            `json:"auto_saved"`
	CreatedAt     time.Time       `json:"created_at"`
}

type EssayTheme struct {
	ID             string          `json:"id"`
	EssayID        string          `json:"essay_id"`
	Colors         json.RawMessage `json:"colors"`
	Fonts          json.RawMessage `json:"fonts"`
	Layout         json.RawMessage `json:"layout"`
	BGType         string          `json:"bg_type"`
	BGGradient     *BGGradient     `json:"bg_gradient,omitempty"`
	BGImageID      *string         `json:"bg_image_id,omitempty"`
	GeneratedCSS   *string         `json:"generated_css,omitempty"`
	CSSGeneratedAt *time.Time      `json:"css_generated_at,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type ServerTheme struct {
	ID             string          `json:"id"`
	Colors         json.RawMessage `json:"colors"`
	Fonts          json.RawMessage `json:"fonts"`
	Layout         json.RawMessage `json:"layout"`
	BGType         string          `json:"bg_type"`
	GeneratedCSS   *string         `json:"generated_css,omitempty"`
	CSSGeneratedAt *time.Time      `json:"css_generated_at,omitempty"`
	UpdatedBy      *string         `json:"updated_by,omitempty"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type PostStyle struct {
	ID                string  `json:"id"`
	AccountID         string  `json:"account_id"`
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
	GeneratedCSS      *string `json:"generated_css,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type Repository interface {
	CreatePageTheme(ctx context.Context, theme *PageTheme) (*PageTheme, error)
	GetPageThemeByID(ctx context.Context, id string) (*PageTheme, error)
	GetPageThemeByAccountID(ctx context.Context, accountID string) (*PageTheme, error)
	UpdatePageTheme(ctx context.Context, theme *PageTheme) error
	DeletePageTheme(ctx context.Context, id string) error

	CreateThemeVersion(ctx context.Context, version *ThemeVersion) (*ThemeVersion, error)
	ListThemeVersions(ctx context.Context, accountID string, limit, offset int) ([]*ThemeVersion, error)
	GetThemeVersionByID(ctx context.Context, id string) (*ThemeVersion, error)
	DeleteThemeVersion(ctx context.Context, id string) error
	CountAutoSavedVersions(ctx context.Context, accountID string) (int, error)
	CountNamedVersions(ctx context.Context, accountID string) (int, error)
	DeleteOldestAutoSave(ctx context.Context, accountID string) error

	CreateEssayTheme(ctx context.Context, theme *EssayTheme) (*EssayTheme, error)
	GetEssayThemeByEssayID(ctx context.Context, essayID string) (*EssayTheme, error)
	UpdateEssayTheme(ctx context.Context, theme *EssayTheme) error
	DeleteEssayTheme(ctx context.Context, essayID string) error

	GetServerTheme(ctx context.Context) (*ServerTheme, error)
	UpsertServerTheme(ctx context.Context, theme *ServerTheme) error

	CreatePostStyle(ctx context.Context, style *PostStyle) (*PostStyle, error)
	GetPostStyleByID(ctx context.Context, id string) (*PostStyle, error)
	ListPostStylesByAccount(ctx context.Context, accountID string, limit, offset int) ([]*PostStyle, error)
	DeletePostStyle(ctx context.Context, id string) error
}
