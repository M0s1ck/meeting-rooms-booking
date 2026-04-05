package create

import "github.com/google/uuid"

type Request struct {
	SlotID               uuid.UUID
	CreateConferenceLink bool
}
