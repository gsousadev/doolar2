package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidSlug = errors.New("invalid slug: must contain only lowercase letters, numbers, and underscores")
	ErrEmptySlug   = errors.New("slug cannot be empty")
)

var slugRegex = regexp.MustCompile(`^[a-z0-9_]+$`)

type SlugValueObject struct {
	value string
}

func NewSlug(value string) (SlugValueObject, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return SlugValueObject{}, ErrEmptySlug
	}

	if !slugRegex.MatchString(value) {
		return SlugValueObject{}, ErrInvalidSlug
	}

	return SlugValueObject{value: value}, nil
}

func NewSlugFromString(s string) (SlugValueObject, error) {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")

	// Remove invalid characters
	reg := regexp.MustCompile(`[^a-z0-9_]`)
	s = reg.ReplaceAllString(s, "")

	// Remove consecutive underscores
	reg = regexp.MustCompile(`_+`)
	s = reg.ReplaceAllString(s, "_")

	// Remove leading/trailing underscores
	s = strings.Trim(s, "_")

	return NewSlug(s)
}

func (s SlugValueObject) Value() string {
	return s.value
}

func (s SlugValueObject) Equals(other SlugValueObject) bool {
	return s.value == other.value
}
