package authjwt

import (
	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

type Identity struct {
	UserID uuid.UUID
	Role   user.Role
}
