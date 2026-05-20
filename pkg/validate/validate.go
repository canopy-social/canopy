package validate

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	// hexColorRe matches a strict 6-digit hex color (e.g. #1a2b3c).
	hexColorRe = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

	// usernameRe matches valid local usernames.
	usernameRe = regexp.MustCompile(`^[a-zA-Z0-9_]{1,30}$`)

	// emailRe is a basic email validation pattern.
	emailRe = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// bcp47LanguageRe matches simple BCP 47 language tags (e.g. en, en-US).
	bcp47LanguageRe = regexp.MustCompile(`^[a-zA-Z]{2,3}(-[a-zA-Z0-9]{2,8})*$`)

	// reservedUsernames cannot be used for registration.
	reservedUsernames = map[string]bool{
		"admin": true, "administrator": true, "root": true,
		"system": true, "moderator": true, "mod": true,
		"support": true, "help": true, "info": true,
		"api": true, "www": true, "mail": true,
		"inbox": true, "outbox": true, "users": true,
		"channels": true, "settings": true, "explore": true,
		"search": true, "login": true, "register": true,
		"auth": true, "oauth": true, "well-known": true,
		"nodeinfo": true, "feed": true, "embed": true,
		"oembed": true, "health": true, "status": true,
	}
)

// HexColor validates a 6-digit hex color string (e.g. #1a2b3c).
func HexColor(s string) bool {
	return hexColorRe.MatchString(s)
}

// Username validates a local username.
func Username(s string) error {
	if !usernameRe.MatchString(s) {
		return fmt.Errorf("username must be 1-30 characters, alphanumeric or underscore")
	}
	if reservedUsernames[strings.ToLower(s)] {
		return fmt.Errorf("username %q is reserved", s)
	}
	return nil
}

// Email validates an email address format.
func Email(s string) bool {
	return emailRe.MatchString(s)
}

// URL validates that a string is a valid HTTP(S) URL.
func URL(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

// Language validates a BCP 47 language tag.
func Language(s string) bool {
	return bcp47LanguageRe.MatchString(s)
}

// FontFamily checks if a font family name is in the allowed list.
func FontFamily(family string, allowed []string) bool {
	for _, f := range allowed {
		if strings.EqualFold(f, family) {
			return true
		}
	}
	return false
}

// IntInRange checks if an integer is within [min, max].
func IntInRange(val, min, max int) bool {
	return val >= min && val <= max
}

// FloatInRange checks if a float is within [min, max].
func FloatInRange(val, min, max float64) bool {
	return val >= min && val <= max
}

// FontWeight validates a CSS font weight (100–900 in steps of 100).
func FontWeight(w int) bool {
	return w >= 100 && w <= 900 && w%100 == 0
}

// Visibility validates a post/essay visibility level.
func Visibility(v string) bool {
	switch v {
	case "public", "unlisted", "followers", "direct":
		return true
	}
	return false
}

// ActorType validates an ActivityPub actor type.
func ActorType(t string) bool {
	switch t {
	case "Person", "Group", "Application", "Service":
		return true
	}
	return false
}
