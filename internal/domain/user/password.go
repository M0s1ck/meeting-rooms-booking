package user

import (
	"errors"
	"unicode"
)

var (
	ErrPassTooShort = errors.New("password must be at least 8 characters")
	ErrPassTooLong  = errors.New("password is too long")
	ErrPassNoLetter = errors.New("password must contain a letter")
	ErrPassNoDigit  = errors.New("password must contain a digit")
)

const (
	passMinLen = 8
	passMaxLen = 50
)

func ValidatePassStrength(password string) error {
	if len([]rune(password)) < passMinLen {
		return ErrPassTooShort
	}

	if len([]rune(password)) > passMaxLen {
		return ErrPassTooLong
	}

	var hasLetter, hasDigit bool

	for _, r := range password {
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
	}

	if !hasLetter {
		return ErrPassNoLetter
	}
	if !hasDigit {
		return ErrPassNoDigit
	}

	return nil
}
