package validate

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	hexColorRe = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

	usernameRe = regexp.MustCompile(`^[a-zA-Z0-9_]{1,30}$`)

	emailRe = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	bcp47LanguageRe = regexp.MustCompile(`^[a-zA-Z]{2,3}(-[a-zA-Z0-9]{2,8})*$`)

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

func HexColor(s string) bool {
	return hexColorRe.MatchString(s)
}

func Username(s string) error {
	if !usernameRe.MatchString(s) {
		return fmt.Errorf("username must be 1-30 characters, alphanumeric or underscore")
	}
	if reservedUsernames[strings.ToLower(s)] {
		return fmt.Errorf("username %q is reserved", s)
	}
	return nil
}

func Email(s string) bool {
	return emailRe.MatchString(s)
}

func URL(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

func Language(s string) bool {
	return bcp47LanguageRe.MatchString(s)
}

func FontFamily(family string, allowed []string) bool {
	for _, f := range allowed {
		if strings.EqualFold(f, family) {
			return true
		}
	}
	return false
}

func IntInRange(val, min, max int) bool {
	return val >= min && val <= max
}

func FloatInRange(val, min, max float64) bool {
	return val >= min && val <= max
}

func FontWeight(w int) bool {
	return w >= 100 && w <= 900 && w%100 == 0
}

func Visibility(v string) bool {
	switch v {
	case "public", "unlisted", "followers", "direct":
		return true
	}
	return false
}

func ActorType(t string) bool {
	switch t {
	case "Person", "Group", "Application", "Service":
		return true
	}
	return false
}
