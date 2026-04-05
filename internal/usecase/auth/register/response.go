package register

import (
	"time"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

type Response struct {
	ID        uuid.UUID
	Email     user.Email
	Role      user.Role
	CreatedAt time.Time
}
