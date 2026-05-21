package themes

import (
	"fmt"
	"strings"
)

func GeneratePageCSS(theme *PageTheme) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("[data-theme-id=%q] {\n", theme.ID))

	writeColorVars(&b, &theme.Colors)
	writeFontVars(&b, &theme.Fonts)
	writeLayoutVars(&b, theme)
	writeBackground(&b, theme)

	b.WriteString("}\n")

	return b.String()
}

func GeneratePostStyleCSS(style *PostStyle) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("[data-post-style=%q] {\n", style.ID))

	if style.BackgroundColor != nil {
		b.WriteString(fmt.Sprintf("  --post-bg: %s;\n", *style.BackgroundColor))
		b.WriteString(fmt.Sprintf("  background-color: %s;\n", *style.BackgroundColor))
	}
	if style.TextColor != nil {
		b.WriteString(fmt.Sprintf("  --post-text: %s;\n", *style.TextColor))
		b.WriteString(fmt.Sprintf("  color: %s;\n", *style.TextColor))
	}
	if style.FontFamily != nil {
		safe := SanitizeFontFamily(*style.FontFamily)
		b.WriteString(fmt.Sprintf("  font-family: %q, system-ui;\n", safe))
	}
	if style.FontSize != nil {
		b.WriteString(fmt.Sprintf("  font-size: %dpx;\n", *style.FontSize))
	}
	if style.FontWeight != nil {
		b.WriteString(fmt.Sprintf("  font-weight: %d;\n", *style.FontWeight))
	}
	if style.BorderRadius != nil {
		b.WriteString(fmt.Sprintf("  border-radius: %dpx;\n", *style.BorderRadius))
	}
	if style.BorderColor != nil {
		b.WriteString(fmt.Sprintf("  border-color: %s;\n", *style.BorderColor))
	}
	if style.BorderWidth != nil {
		b.WriteString(fmt.Sprintf("  border-width: %dpx;\n", *style.BorderWidth))
		b.WriteString("  border-style: solid;\n")
	}
	if style.Padding != nil {
		b.WriteString(fmt.Sprintf("  padding: %dpx;\n", *style.Padding))
	}

	b.WriteString("}\n")

	return b.String()
}

func GenerateEssayThemeCSS(theme *EssayTheme, base *PageTheme) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("[data-essay-theme=%q] {\n", theme.ID))

	if theme.BGType != "inherit" && theme.BGType != "" {
		b.WriteString(fmt.Sprintf("  --essay-bg-type: %s;\n", theme.BGType))
	}

	b.WriteString("}\n")

	return b.String()
}

func GenerateServerThemeCSS(theme *ServerTheme) string {
	var b strings.Builder

	b.WriteString(":root[data-server-theme] {\n")
	b.WriteString("}\n")

	return b.String()
}

func writeColorVars(b *strings.Builder, c *ColorPalette) {
	b.WriteString(fmt.Sprintf("  --color-background: %s;\n", c.Background))
	b.WriteString(fmt.Sprintf("  --color-surface: %s;\n", c.Surface))
	b.WriteString(fmt.Sprintf("  --color-text-primary: %s;\n", c.TextPrimary))
	b.WriteString(fmt.Sprintf("  --color-text-secondary: %s;\n", c.TextSecondary))
	b.WriteString(fmt.Sprintf("  --color-accent: %s;\n", c.Accent))
	b.WriteString(fmt.Sprintf("  --color-accent-text: %s;\n", c.AccentText))
	b.WriteString(fmt.Sprintf("  --color-border: %s;\n", c.Border))
	b.WriteString(fmt.Sprintf("  --color-link: %s;\n", c.Link))
}

func writeFontVars(b *strings.Builder, f *FontPalette) {
	writeSingleFont(b, "body", &f.Body)
	writeSingleFont(b, "heading", &f.Heading)
	writeSingleFont(b, "mono", &f.Mono)
	writeSingleFont(b, "display", &f.Display)
}

func writeSingleFont(b *strings.Builder, prefix string, fc *FontConfig) {
	safe := SanitizeFontFamily(fc.Family)
	b.WriteString(fmt.Sprintf("  --font-%s-family: %q, system-ui;\n", prefix, safe))
	b.WriteString(fmt.Sprintf("  --font-%s-size: %dpx;\n", prefix, fc.Size))
	b.WriteString(fmt.Sprintf("  --font-%s-weight: %d;\n", prefix, fc.Weight))
	b.WriteString(fmt.Sprintf("  --font-%s-line-height: %.2f;\n", prefix, fc.LineHeight))
	b.WriteString(fmt.Sprintf("  --font-%s-letter-spacing: %.3fem;\n", prefix, fc.LetterSpacing))
}

func writeLayoutVars(b *strings.Builder, t *PageTheme) {
	b.WriteString(fmt.Sprintf("  --page-max-width: %dpx;\n", t.PageMaxWidth))
	b.WriteString(fmt.Sprintf("  --page-padding: %dpx;\n", t.PagePadding))
	b.WriteString(fmt.Sprintf("  --bg-blur: %dpx;\n", t.BGBlur))
	b.WriteString(fmt.Sprintf("  --bg-opacity: %d%%;\n", t.BGOpacity))
}

func writeBackground(b *strings.Builder, t *PageTheme) {
	switch t.BGType {
	case "color":
		b.WriteString(fmt.Sprintf("  background-color: %s;\n", t.Colors.Background))
	case "gradient":
		if t.BGGradient != nil {
			b.WriteString(fmt.Sprintf("  background: %s;\n", buildGradientCSS(t.BGGradient)))
		}
	case "image":
		if t.BGImageID != nil {
			b.WriteString(fmt.Sprintf("  background-image: url('/api/v1/media/%s');\n", *t.BGImageID))
			b.WriteString(fmt.Sprintf("  background-size: %s;\n", t.BGImageSize))
			b.WriteString("  background-position: center;\n")
			b.WriteString("  background-repeat: no-repeat;\n")
		}
	}
}

func buildGradientCSS(g *BGGradient) string {
	stops := make([]string, len(g.Stops))
	for i, s := range g.Stops {
		stops[i] = fmt.Sprintf("%s %.1f%%", s.Color, s.Position)
	}

	switch g.Type {
	case "radial":
		return fmt.Sprintf("radial-gradient(circle, %s)", strings.Join(stops, ", "))
	case "conic":
		return fmt.Sprintf("conic-gradient(from %ddeg, %s)", g.Angle, strings.Join(stops, ", "))
	default:
		return fmt.Sprintf("linear-gradient(%ddeg, %s)", g.Angle, strings.Join(stops, ", "))
	}
}
