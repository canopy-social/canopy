package themes

import (
	"fmt"
	"regexp"
	"strings"
)

var hexColorRegex = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)

var allowedFonts = map[string]bool{
	"system-ui":       true,
	"Inter":           true,
	"Outfit":          true,
	"Roboto":          true,
	"Roboto Mono":     true,
	"Open Sans":       true,
	"Lato":            true,
	"Montserrat":      true,
	"Poppins":         true,
	"Raleway":         true,
	"Playfair Display": true,
	"Merriweather":    true,
	"Source Code Pro":  true,
	"Fira Code":       true,
	"JetBrains Mono":  true,
	"IBM Plex Sans":   true,
	"IBM Plex Mono":   true,
	"DM Sans":         true,
	"DM Serif Display": true,
	"Space Grotesk":   true,
	"Space Mono":      true,
	"Archivo":         true,
	"Sora":            true,
	"Inconsolata":     true,
	"monospace":       true,
	"serif":           true,
	"sans-serif":      true,
	"cursive":         true,
}

var allowedBGTypes = map[string]bool{
	"color":    true,
	"gradient": true,
	"image":    true,
	"inherit":  true,
}

var allowedBGSizes = map[string]bool{
	"cover":   true,
	"contain": true,
	"auto":    true,
}

var allowedTextureTypes = map[string]bool{
	"noise":     true,
	"grain":     true,
	"dots":      true,
	"lines":     true,
	"crosshatch": true,
}

func ValidateHexColor(color string) error {
	if !hexColorRegex.MatchString(color) {
		return fmt.Errorf("invalid hex color: %s", color)
	}
	return nil
}

func ValidateColorPalette(c *ColorPalette) error {
	colors := map[string]string{
		"background":     c.Background,
		"surface":        c.Surface,
		"text_primary":   c.TextPrimary,
		"text_secondary": c.TextSecondary,
		"accent":         c.Accent,
		"accent_text":    c.AccentText,
		"border":         c.Border,
		"link":           c.Link,
	}
	for name, val := range colors {
		if err := ValidateHexColor(val); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	return nil
}

func ValidateFontConfig(name string, f *FontConfig) error {
	if !allowedFonts[f.Family] {
		return fmt.Errorf("%s: font family %q is not allowed", name, f.Family)
	}
	if f.Size < 8 || f.Size > 128 {
		return fmt.Errorf("%s: font size must be between 8 and 128", name)
	}
	if f.Weight < 100 || f.Weight > 900 || f.Weight%100 != 0 {
		return fmt.Errorf("%s: font weight must be a multiple of 100 between 100-900", name)
	}
	if f.LineHeight < 0.5 || f.LineHeight > 3.0 {
		return fmt.Errorf("%s: line height must be between 0.5 and 3.0", name)
	}
	if f.LetterSpacing < -0.1 || f.LetterSpacing > 0.5 {
		return fmt.Errorf("%s: letter spacing must be between -0.1 and 0.5", name)
	}
	return nil
}

func ValidateFontPalette(p *FontPalette) error {
	validators := []struct {
		name string
		cfg  *FontConfig
	}{
		{"body", &p.Body},
		{"heading", &p.Heading},
		{"mono", &p.Mono},
		{"display", &p.Display},
	}
	for _, v := range validators {
		if err := ValidateFontConfig(v.name, v.cfg); err != nil {
			return err
		}
	}
	return nil
}

func ValidateBGType(bgType string) error {
	if !allowedBGTypes[bgType] {
		return fmt.Errorf("invalid background type: %s", bgType)
	}
	return nil
}

func ValidateBGSize(bgSize string) error {
	if !allowedBGSizes[bgSize] {
		return fmt.Errorf("invalid background size: %s", bgSize)
	}
	return nil
}

func ValidateGradient(g *BGGradient) error {
	if g == nil {
		return nil
	}
	validTypes := map[string]bool{"linear": true, "radial": true, "conic": true}
	if !validTypes[g.Type] {
		return fmt.Errorf("invalid gradient type: %s", g.Type)
	}
	if g.Angle < 0 || g.Angle > 360 {
		return fmt.Errorf("gradient angle must be between 0 and 360")
	}
	if len(g.Stops) < 2 || len(g.Stops) > 10 {
		return fmt.Errorf("gradient must have between 2 and 10 stops")
	}
	for i, stop := range g.Stops {
		if err := ValidateHexColor(stop.Color); err != nil {
			return fmt.Errorf("gradient stop %d: %w", i, err)
		}
		if stop.Position < 0 || stop.Position > 100 {
			return fmt.Errorf("gradient stop %d: position must be between 0 and 100", i)
		}
	}
	return nil
}

func ValidatePageTheme(t *PageTheme) error {
	if err := ValidateColorPalette(&t.Colors); err != nil {
		return fmt.Errorf("colors: %w", err)
	}
	if err := ValidateFontPalette(&t.Fonts); err != nil {
		return fmt.Errorf("fonts: %w", err)
	}
	if err := ValidateBGType(t.BGType); err != nil {
		return err
	}
	if err := ValidateBGSize(t.BGImageSize); err != nil {
		return err
	}
	if t.BGType == "gradient" {
		if err := ValidateGradient(t.BGGradient); err != nil {
			return err
		}
	}
	if t.BGBlur < 0 || t.BGBlur > 50 {
		return fmt.Errorf("background blur must be between 0 and 50")
	}
	if t.BGOpacity < 0 || t.BGOpacity > 100 {
		return fmt.Errorf("background opacity must be between 0 and 100")
	}
	if t.PageMaxWidth < 400 || t.PageMaxWidth > 1600 {
		return fmt.Errorf("page max width must be between 400 and 1600")
	}
	if t.PagePadding < 0 || t.PagePadding > 64 {
		return fmt.Errorf("page padding must be between 0 and 64")
	}
	if len(t.Stickers) > 50 {
		return fmt.Errorf("maximum 50 stickers allowed")
	}
	if len(t.Layout) > 20 {
		return fmt.Errorf("maximum 20 layout widgets allowed")
	}
	return nil
}

func ValidatePostStyle(s *PostStyle) error {
	if s.BackgroundColor != nil {
		if err := ValidateHexColor(*s.BackgroundColor); err != nil {
			return fmt.Errorf("background_color: %w", err)
		}
	}
	if s.TextColor != nil {
		if err := ValidateHexColor(*s.TextColor); err != nil {
			return fmt.Errorf("text_color: %w", err)
		}
	}
	if s.BorderColor != nil {
		if err := ValidateHexColor(*s.BorderColor); err != nil {
			return fmt.Errorf("border_color: %w", err)
		}
	}
	if s.FontFamily != nil && !allowedFonts[*s.FontFamily] {
		return fmt.Errorf("font family %q is not allowed", *s.FontFamily)
	}
	if s.FontSize != nil && (*s.FontSize < 8 || *s.FontSize > 72) {
		return fmt.Errorf("font size must be between 8 and 72")
	}
	if s.FontWeight != nil && (*s.FontWeight < 100 || *s.FontWeight > 900 || *s.FontWeight%100 != 0) {
		return fmt.Errorf("font weight must be a multiple of 100 between 100-900")
	}
	if s.BorderRadius != nil && (*s.BorderRadius < 0 || *s.BorderRadius > 50) {
		return fmt.Errorf("border radius must be between 0 and 50")
	}
	if s.BorderWidth != nil && (*s.BorderWidth < 0 || *s.BorderWidth > 10) {
		return fmt.Errorf("border width must be between 0 and 10")
	}
	if s.Padding != nil && (*s.Padding < 0 || *s.Padding > 64) {
		return fmt.Errorf("padding must be between 0 and 64")
	}
	if s.TextureType != nil && !allowedTextureTypes[*s.TextureType] {
		return fmt.Errorf("invalid texture type: %s", *s.TextureType)
	}
	return nil
}

func IsAllowedFont(family string) bool {
	return allowedFonts[family]
}

func ListAllowedFonts() []string {
	fonts := make([]string, 0, len(allowedFonts))
	for f := range allowedFonts {
		fonts = append(fonts, f)
	}
	return fonts
}

func SanitizeFontFamily(family string) string {
	if strings.ContainsAny(family, ";{}()") {
		return "system-ui"
	}
	return family
}
