package authjwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/config"
	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

type Manager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewManager(cfg *config.JwtCfg) *Manager {
	return &Manager{
		secret: []byte(cfg.Secret),
		issuer: cfg.Issuer,
		ttl:    cfg.Ttl,
	}
}

var ErrInvalidToken = errors.New("invalid jwt")

func (m *Manager) Generate(userID uuid.UUID, role user.Role) (string, error) {
	claims := &Claims{
		Role:   role,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    m.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	str, err := token.SignedString(m.secret)
	if err != nil {
		return "", err
	}

	return str, nil
}

func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			return m.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (m *Manager) ParseToIdentity(tokenStr string) (*Identity, error) {
	claims, err := m.Parse(tokenStr)
	if err != nil {
		return nil, err
	}

	return &Identity{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}
