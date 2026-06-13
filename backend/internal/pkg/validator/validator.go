package validator

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

var usernamePattern = regexp.MustCompile(`^[a-z0-9-]+$`)

func Required(value string, field string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(field + " is required")
	}
	return nil
}

func Email(value string) error {
	if !strings.Contains(value, "@") {
		return errors.New("email must be valid")
	}
	return nil
}

func URL(value string) error {
	parsed, err := url.ParseRequestURI(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return errors.New("url must be valid")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("url must use http or https")
	}
	return nil
}

func Username(value string) error {
	if !usernamePattern.MatchString(value) {
		return errors.New("username may only contain lowercase letters, digits, and hyphens")
	}
	return nil
}

func Password(value string) error {
	if len(strings.TrimSpace(value)) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}
