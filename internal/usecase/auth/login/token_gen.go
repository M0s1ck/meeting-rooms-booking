package login

import (
	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

//go:generate mockgen -source=token_gen.go -destination=mocks/token_gen_mock.go -package=mocks
type tokenGenerator interface {
	Generate(userID uuid.UUID, role user.Role) (string, error)
}
