package dummylogin

import (
	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

type TokenGenerator interface {
	Generate(userID uuid.UUID, role user.Role) (string, error)
}
