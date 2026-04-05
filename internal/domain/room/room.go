package room

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Capacity    *int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func New(name string, description *string, capacity *int, now time.Time) (*Room, error) {
	normalizedName := strings.TrimSpace(name)
	if normalizedName == "" {
		return nil, ErrEmptyName
	}

	normalizedDescription := normalizeDescription(description)
	normalizedCapacity, err := normalizeCapacity(capacity)
	if err != nil {
		return nil, err
	}

	timestamp := now.UTC()

	return &Room{
		ID:          uuid.New(),
		Name:        normalizedName,
		Description: normalizedDescription,
		Capacity:    normalizedCapacity,
		CreatedAt:   timestamp,
		UpdatedAt:   timestamp,
	}, nil
}

func normalizeDescription(description *string) *string {
	if description == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*description)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func normalizeCapacity(capacity *int) (*int, error) {
	if capacity == nil {
		return nil, nil
	}

	if *capacity <= 0 {
		return nil, ErrInvalidCapacity
	}

	normalized := *capacity
	return &normalized, nil
}
