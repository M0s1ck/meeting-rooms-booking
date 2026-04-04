package bcrypt

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/internships-backend/test-backend-M0s1ck/internal/config"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

type Hasher struct {
	cost int
}

func NewHasher(cfg *config.BcryptCfg) *Hasher {
	cost := cfg.Cost
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	return &Hasher{cost: cost}
}

func (b *Hasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", fmt.Errorf("pass hash err: %w", err)
	}

	return string(hash), nil
}

func (b *Hasher) Compare(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return user.ErrWrongPassword
	}

	return fmt.Errorf("pass hash compare err: %w", err)
}
