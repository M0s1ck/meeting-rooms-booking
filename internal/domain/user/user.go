package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     Email
	PassHash  string
	Role      Role
	CreatedAt time.Time
}

func New(id uuid.UUID, email Email, passHash string, role Role, now time.Time) (*User, error) {
	if passHash == "" {
		return nil, ErrEmptyPassHash
	}

	return &User{
		ID:        id,
		Email:     email,
		PassHash:  passHash,
		Role:      role,
		CreatedAt: now.UTC(),
	}, nil
}
