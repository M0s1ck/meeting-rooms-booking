package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/internships-backend/test-backend-M0s1ck/internal/service/authjwt"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/helpers"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"
)

type AuthMiddleware struct {
	identityParser identityParser
	logger         *slog.Logger
}

type identityParser interface {
	ParseToIdentity(tokenStr string) (*authjwt.Identity, error)
}

type identityCtxKeyT struct{}

var identityCtxKey = identityCtxKeyT{}

func IdentityFromContext(ctx context.Context) (*authjwt.Identity, bool) {
	identity, ok := ctx.Value(identityCtxKey).(*authjwt.Identity)
	return identity, ok
}

func MustIdentityFromContext(ctx context.Context) *authjwt.Identity {
	identity, ok := IdentityFromContext(ctx)
	if !ok {
		panic("auth identity is missing from context")
	}

	return identity
}

func NewAuthMiddleware(identityParser identityParser, logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		identityParser: identityParser,
		logger:         logger,
	}
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tokenStr, err := extractBearerToken(r.Header.Get("Authorization"))
		if err != nil {
			writeUnauthorized(ctx, m.logger, w, err.Error())
			return
		}

		identity, err := m.identityParser.ParseToIdentity(tokenStr)
		if err != nil {
			writeUnauthorized(ctx, m.logger, w, "invalid or expired bearer token")
			return
		}

		ctx = context.WithValue(ctx, identityCtxKey, identity)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errMissingAuthorizationHeader
	}

	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", errInvalidAuthorizationHeader
	}

	return parts[1], nil
}

func writeUnauthorized(ctx context.Context, logger *slog.Logger, w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	resp := helpers.NewErrorResponse(oapi.UNAUTHORIZED, message)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.ErrorContext(ctx, "couldn't write auth error response", "err", err)
	}
}

var (
	errMissingAuthorizationHeader = errors.New("request should contain Authorization header")
	errInvalidAuthorizationHeader = errors.New("authorization header must be in format: Bearer <token>")
)
