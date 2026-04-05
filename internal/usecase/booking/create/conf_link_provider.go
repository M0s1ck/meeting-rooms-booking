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

//go:generate mockgen -source=conf_link_provider.go -destination=mocks/conf_link_provider_mock.go -package=mocks
type conferenceLinkProvider interface {
	CreateLink(ctx context.Context, slotID uuid.UUID, userID uuid.UUID) (string, error)
}
