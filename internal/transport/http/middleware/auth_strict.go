package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"
)

var publicOperationIDs = map[string]struct{}{
	"PostDummyLogin": {},
	"PostLogin":      {},
	"PostRegister":   {},
}

func NewStrictMiddlewares(identityParser identityParser, logger *slog.Logger) []oapi.StrictMiddlewareFunc {
	return []oapi.StrictMiddlewareFunc{
		StrictAuthMiddleware(identityParser, logger),
	}
}

func StrictAuthMiddleware(identityParser identityParser, logger *slog.Logger) oapi.StrictMiddlewareFunc {
	return func(next oapi.StrictHandlerFunc, operationID string) oapi.StrictHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {

			if _, isPublic := publicOperationIDs[operationID]; isPublic {
				return next(ctx, w, r, request)
			}

			tokenStr, err := extractBearerToken(r.Header.Get("Authorization"))
			if err != nil {
				writeUnauthorized(ctx, logger, w, err.Error())
				return nil, nil
			}

			identity, err := identityParser.ParseToIdentity(tokenStr)
			if err != nil {
				writeUnauthorized(ctx, logger, w, "invalid or expired bearer token")
				return nil, nil
			}

			ctx = context.WithValue(ctx, identityCtxKey, identity)
			return next(ctx, w, r.WithContext(ctx), request)
		}
	}
}
