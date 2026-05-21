package themes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

type mockRepo struct {
	pageThemes    map[string]*PageTheme
	themeVersions map[string]*ThemeVersion
	essayThemes   map[string]*EssayTheme
	serverTheme   *ServerTheme
	postStyles    map[string]*PostStyle
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		pageThemes:    make(map[string]*PageTheme),
		themeVersions: make(map[string]*ThemeVersion),
		essayThemes:   make(map[string]*EssayTheme),
		postStyles:    make(map[string]*PostStyle),
	}
}

func (m *mockRepo) CreatePageTheme(_ context.Context, t *PageTheme) (*PageTheme, error) {
	m.pageThemes[t.ID] = t
	return t, nil
}

func (m *mockRepo) GetPageThemeByID(_ context.Context, id string) (*PageTheme, error) {
	t, ok := m.pageThemes[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return t, nil
}

func (m *mockRepo) GetPageThemeByAccountID(_ context.Context, accountID string) (*PageTheme, error) {
	for _, t := range m.pageThemes {
		if t.AccountID == accountID {
			return t, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockRepo) UpdatePageTheme(_ context.Context, t *PageTheme) error {
	m.pageThemes[t.ID] = t
	return nil
}

func (m *mockRepo) DeletePageTheme(_ context.Context, id string) error {
	delete(m.pageThemes, id)
	return nil
}

func (m *mockRepo) CreateThemeVersion(_ context.Context, v *ThemeVersion) (*ThemeVersion, error) {
	m.themeVersions[v.ID] = v
	return v, nil
}

func (m *mockRepo) ListThemeVersions(_ context.Context, accountID string, limit, offset int) ([]*ThemeVersion, error) {
	var result []*ThemeVersion
	for _, v := range m.themeVersions {
		if v.AccountID == accountID {
			result = append(result, v)
		}
	}
	return result, nil
}

func (m *mockRepo) GetThemeVersionByID(_ context.Context, id string) (*ThemeVersion, error) {
	v, ok := m.themeVersions[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (m *mockRepo) DeleteThemeVersion(_ context.Context, id string) error {
	delete(m.themeVersions, id)
	return nil
}

func (m *mockRepo) CountAutoSavedVersions(_ context.Context, accountID string) (int, error) {
	count := 0
	for _, v := range m.themeVersions {
		if v.AccountID == accountID && v.AutoSaved {
			count++
		}
	}
	return count, nil
}

func (m *mockRepo) CountNamedVersions(_ context.Context, accountID string) (int, error) {
	count := 0
	for _, v := range m.themeVersions {
		if v.AccountID == accountID && !v.AutoSaved {
			count++
		}
	}
	return count, nil
}

func (m *mockRepo) DeleteOldestAutoSave(_ context.Context, accountID string) error {
	var oldest *ThemeVersion
	for _, v := range m.themeVersions {
		if v.AccountID == accountID && v.AutoSaved {
			if oldest == nil || v.CreatedAt.Before(oldest.CreatedAt) {
				oldest = v
			}
		}
	}
	if oldest != nil {
		delete(m.themeVersions, oldest.ID)
	}
	return nil
}

func (m *mockRepo) CreateEssayTheme(_ context.Context, t *EssayTheme) (*EssayTheme, error) {
	m.essayThemes[t.ID] = t
	return t, nil
}

func (m *mockRepo) GetEssayThemeByEssayID(_ context.Context, essayID string) (*EssayTheme, error) {
	for _, t := range m.essayThemes {
		if t.EssayID == essayID {
			return t, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockRepo) UpdateEssayTheme(_ context.Context, t *EssayTheme) error {
	for k, v := range m.essayThemes {
		if v.EssayID == t.EssayID {
			m.essayThemes[k] = t
			return nil
		}
	}
	return errors.New("not found")
}

func (m *mockRepo) DeleteEssayTheme(_ context.Context, essayID string) error {
	for k, v := range m.essayThemes {
		if v.EssayID == essayID {
			delete(m.essayThemes, k)
			return nil
		}
	}
	return nil
}

func (m *mockRepo) GetServerTheme(_ context.Context) (*ServerTheme, error) {
	if m.serverTheme == nil {
		return nil, errors.New("not found")
	}
	return m.serverTheme, nil
}

func (m *mockRepo) UpsertServerTheme(_ context.Context, t *ServerTheme) error {
	m.serverTheme = t
	return nil
}

func (m *mockRepo) CreatePostStyle(_ context.Context, s *PostStyle) (*PostStyle, error) {
	m.postStyles[s.ID] = s
	return s, nil
}

func (m *mockRepo) GetPostStyleByID(_ context.Context, id string) (*PostStyle, error) {
	s, ok := m.postStyles[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return s, nil
}

func (m *mockRepo) ListPostStylesByAccount(_ context.Context, accountID string, limit, offset int) ([]*PostStyle, error) {
	var result []*PostStyle
	for _, s := range m.postStyles {
		if s.AccountID == accountID {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *mockRepo) DeletePostStyle(_ context.Context, id string) error {
	delete(m.postStyles, id)
	return nil
}

func TestCreatePageTheme(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	theme, err := svc.CreatePageTheme(context.Background(), "acct_001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if theme.AccountID != "acct_001" {
		t.Errorf("expected account_id acct_001, got %s", theme.AccountID)
	}
	if theme.Colors.Background != "#ffffff" {
		t.Errorf("expected default background #ffffff, got %s", theme.Colors.Background)
	}
	if theme.GeneratedCSS == nil {
		t.Error("expected generated CSS to be set")
	}
}

func TestCreatePageThemeDuplicate(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, _ = svc.CreatePageTheme(context.Background(), "acct_001")
	_, err := svc.CreatePageTheme(context.Background(), "acct_001")
	if err == nil {
		t.Error("expected duplicate error")
	}
}

func TestUpdatePageTheme(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, _ = svc.CreatePageTheme(context.Background(), "acct_001")

	accent := "#ff0066"
	updated, err := svc.UpdatePageTheme(context.Background(), "acct_001", &PageThemeUpdateRequest{
		Colors: &ColorPalette{
			Background:    "#000000",
			Surface:       "#111111",
			TextPrimary:   "#ffffff",
			TextSecondary: "#888888",
			Accent:        accent,
			AccentText:    "#ffffff",
			Border:        "#333333",
			Link:          "#ff0066",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Colors.Accent != accent {
		t.Errorf("expected accent %s, got %s", accent, updated.Colors.Accent)
	}
	if len(repo.themeVersions) == 0 {
		t.Error("expected auto-save version to be created")
	}
}

func TestUpdatePageThemeInvalidColor(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, _ = svc.CreatePageTheme(context.Background(), "acct_001")

	_, err := svc.UpdatePageTheme(context.Background(), "acct_001", &PageThemeUpdateRequest{
		Colors: &ColorPalette{
			Background:    "not-a-color",
			Surface:       "#111111",
			TextPrimary:   "#ffffff",
			TextSecondary: "#888888",
			Accent:        "#ff0066",
			AccentText:    "#ffffff",
			Border:        "#333333",
			Link:          "#ff0066",
		},
	})
	if err == nil {
		t.Error("expected validation error for invalid color")
	}
}

func TestSaveNamedVersion(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, _ = svc.CreatePageTheme(context.Background(), "acct_001")

	version, err := svc.SaveNamedVersion(context.Background(), "acct_001", "My Dark Theme")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version.Label == nil || *version.Label != "My Dark Theme" {
		t.Error("expected label to be 'My Dark Theme'")
	}
	if version.AutoSaved {
		t.Error("expected auto_saved to be false")
	}
}

func TestRestoreVersion(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, _ = svc.CreatePageTheme(context.Background(), "acct_001")
	version, _ := svc.SaveNamedVersion(context.Background(), "acct_001", "Snapshot")

	accent := "#00ff00"
	_, _ = svc.UpdatePageTheme(context.Background(), "acct_001", &PageThemeUpdateRequest{
		Colors: &ColorPalette{
			Background:    "#000000",
			Surface:       "#111111",
			TextPrimary:   "#ffffff",
			TextSecondary: "#888888",
			Accent:        accent,
			AccentText:    "#ffffff",
			Border:        "#333333",
			Link:          accent,
		},
	})

	restored, err := svc.RestoreVersion(context.Background(), "acct_001", version.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if restored.Colors.Accent == accent {
		t.Error("expected restored theme to not have the updated accent")
	}
}

func TestCreatePostStyle(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	bg := "#ff0000"
	style, err := svc.CreatePostStyle(context.Background(), "acct_001", &CreatePostStyleRequest{
		BackgroundColor: &bg,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if style.BackgroundColor == nil || *style.BackgroundColor != bg {
		t.Error("expected background color to match")
	}
	if style.GeneratedCSS == nil {
		t.Error("expected generated CSS to be set")
	}
}

func TestCreatePostStyleInvalidColor(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	bg := "not-valid"
	_, err := svc.CreatePostStyle(context.Background(), "acct_001", &CreatePostStyleRequest{
		BackgroundColor: &bg,
	})
	if err == nil {
		t.Error("expected validation error for invalid color")
	}
}

func TestDeletePostStyleOwnership(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	bg := "#ff0000"
	style, _ := svc.CreatePostStyle(context.Background(), "acct_001", &CreatePostStyleRequest{
		BackgroundColor: &bg,
	})

	err := svc.DeletePostStyle(context.Background(), "acct_999", style.ID)
	if err == nil {
		t.Error("expected ownership error")
	}
}

func TestEssayThemeCreateAndUpdate(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	theme, err := svc.GetOrCreateEssayTheme(context.Background(), "essay_001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if theme.EssayID != "essay_001" {
		t.Error("expected essay_id to match")
	}
	if theme.BGType != "inherit" {
		t.Errorf("expected bg_type 'inherit', got %s", theme.BGType)
	}

	updated, err := svc.UpdateEssayTheme(context.Background(), "essay_001", json.RawMessage(`{"accent":"#ff0000"}`), nil, nil, "color")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.BGType != "color" {
		t.Errorf("expected bg_type 'color', got %s", updated.BGType)
	}
}

func TestServerTheme(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	theme, err := svc.UpdateServerTheme(context.Background(), "admin_001", json.RawMessage(`{"accent":"#0066ff"}`), nil, nil, "color")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if theme.UpdatedBy == nil || *theme.UpdatedBy != "admin_001" {
		t.Error("expected updated_by to be admin_001")
	}

	fetched, err := svc.GetServerTheme(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fetched.ID != "singleton" {
		t.Errorf("expected singleton id, got %s", fetched.ID)
	}
}

func TestValidateHexColor(t *testing.T) {
	valid := []string{"#fff", "#000000", "#ff0066", "#FF0066", "#aabbcc", "#aabbccdd"}
	for _, c := range valid {
		if err := ValidateHexColor(c); err != nil {
			t.Errorf("expected %s to be valid, got error: %v", c, err)
		}
	}

	invalid := []string{"fff", "#gg0000", "#12345", "rgb(0,0,0)", ""}
	for _, c := range invalid {
		if err := ValidateHexColor(c); err == nil {
			t.Errorf("expected %s to be invalid", c)
		}
	}
}

func TestValidateFontConfig(t *testing.T) {
	valid := FontConfig{Family: "Inter", Size: 16, Weight: 400, LineHeight: 1.6, LetterSpacing: 0}
	if err := ValidateFontConfig("body", &valid); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	disallowed := FontConfig{Family: "Comic Sans", Size: 16, Weight: 400, LineHeight: 1.6, LetterSpacing: 0}
	if err := ValidateFontConfig("body", &disallowed); err == nil {
		t.Error("expected error for disallowed font family")
	}

	tooSmall := FontConfig{Family: "Inter", Size: 4, Weight: 400, LineHeight: 1.6, LetterSpacing: 0}
	if err := ValidateFontConfig("body", &tooSmall); err == nil {
		t.Error("expected error for too small font size")
	}
}

func TestValidateGradient(t *testing.T) {
	valid := &BGGradient{
		Type:  "linear",
		Angle: 45,
		Stops: []GradientStop{
			{Color: "#000000", Position: 0},
			{Color: "#ffffff", Position: 100},
		},
	}
	if err := ValidateGradient(valid); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	tooFewStops := &BGGradient{
		Type:  "linear",
		Angle: 0,
		Stops: []GradientStop{{Color: "#000000", Position: 0}},
	}
	if err := ValidateGradient(tooFewStops); err == nil {
		t.Error("expected error for too few stops")
	}
}

func TestGeneratePageCSS(t *testing.T) {
	theme := &PageTheme{
		ID:        "theme_001",
		AccountID: "acct_001",
		Colors:    DefaultColorPalette(),
		Fonts:     DefaultFontPalette(),
		BGType:    "color",
		BGImageSize: "cover",
		BGBlur:      0,
		BGOpacity:   100,
		PageMaxWidth: 800,
		PagePadding:  24,
	}

	css := GeneratePageCSS(theme)
	if !strings.Contains(css, "--color-background") {
		t.Error("expected CSS to contain --color-background")
	}
	if !strings.Contains(css, "--font-body-family") {
		t.Error("expected CSS to contain --font-body-family")
	}
	if !strings.Contains(css, "--page-max-width") {
		t.Error("expected CSS to contain --page-max-width")
	}
	if !strings.Contains(css, "background-color:") {
		t.Error("expected CSS to contain background-color")
	}
}

func TestGeneratePostStyleCSS(t *testing.T) {
	bg := "#ff0000"
	txt := "#ffffff"
	style := &PostStyle{
		ID:              "ps_001",
		AccountID:       "acct_001",
		BackgroundColor: &bg,
		TextColor:       &txt,
	}

	css := GeneratePostStyleCSS(style)
	if !strings.Contains(css, "--post-bg: #ff0000") {
		t.Error("expected CSS to contain --post-bg")
	}
	if !strings.Contains(css, "--post-text: #ffffff") {
		t.Error("expected CSS to contain --post-text")
	}
}

func TestAutoSaveRotation(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, _ = svc.CreatePageTheme(context.Background(), "acct_001")

	for i := 0; i < 52; i++ {
		repo.themeVersions["v_"+string(rune(i))] = &ThemeVersion{
			ID:        "v_" + string(rune(i)),
			AccountID: "acct_001",
			AutoSaved: true,
			CreatedAt: time.Now().Add(time.Duration(-i) * time.Minute),
		}
	}

	_, err := svc.UpdatePageTheme(context.Background(), "acct_001", &PageThemeUpdateRequest{
		BGBlur: intPtr(5),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNamedVersionLimit(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	_, _ = svc.CreatePageTheme(context.Background(), "acct_001")

	for i := 0; i < 20; i++ {
		label := "Version"
		repo.themeVersions["nv_"+string(rune(i))] = &ThemeVersion{
			ID:        "nv_" + string(rune(i)),
			AccountID: "acct_001",
			AutoSaved: false,
			Label:     &label,
		}
	}

	_, err := svc.SaveNamedVersion(context.Background(), "acct_001", "21st version")
	if err == nil {
		t.Error("expected error for exceeding named version limit")
	}
}

func TestHandlerListAllowedFonts(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	req := httptest.NewRequest("GET", "/api/v1/fonts", nil)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/api/v1/fonts", h.ListAllowedFonts)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var body map[string]interface{}
	json.NewDecoder(w.Body).Decode(&body)
	fonts, ok := body["fonts"].([]interface{})
	if !ok || len(fonts) == 0 {
		t.Error("expected non-empty fonts list")
	}
}

func TestSanitizeFontFamily(t *testing.T) {
	safe := SanitizeFontFamily("Inter")
	if safe != "Inter" {
		t.Errorf("expected Inter, got %s", safe)
	}

	malicious := SanitizeFontFamily("Inter; body { display: none }")
	if malicious != "system-ui" {
		t.Errorf("expected system-ui for malicious input, got %s", malicious)
	}
}

func intPtr(i int) *int {
	return &i
}
