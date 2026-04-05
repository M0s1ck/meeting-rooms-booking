package create

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrConferenceServiceUnavailable = errors.New("conference service unavailable")
	ErrConferenceInternal           = errors.New("conference internal error")
)

type conferenceLinkProvider interface {
	CreateLink(ctx context.Context, slotID uuid.UUID, userID uuid.UUID) (string, error)
}
