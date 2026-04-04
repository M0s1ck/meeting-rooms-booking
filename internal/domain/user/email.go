package user

import (
	"errors"
	"strings"
)

type Email string

var (
	ErrInvalidEmail = errors.New("invalid email format")
	ErrBlankEmail   = errors.New("email is blank")
)

func NewEmail(v string) (Email, error) {
	v = strings.TrimSpace(strings.ToLower(v))

	if v == "" {
		return "", ErrBlankEmail
	}

	if !strings.Contains(v, "@") {
		return "", ErrInvalidEmail
	}

	email := Email(v)
	return email, nil
}
