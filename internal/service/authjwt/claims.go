package authjwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

type Claims struct {
	Role   user.Role `json:"role"`
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}
