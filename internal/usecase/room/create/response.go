package create

import (
	"time"

	"github.com/google/uuid"
)

type Response struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Capacity    *int
	CreatedAt   time.Time
}
