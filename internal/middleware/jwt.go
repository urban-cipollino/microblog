package middleware

import (
	"context"
	"net/http"
	"strings"

	"micro-blog/internal/exception"
	"micro-blog/internal/logger"
	"micro-blog/internal/service/auth"
)

func Auth(jwtService *auth.JWTService, logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return HandleWrapper(func(w http.ResponseWriter, r *http.Request) error {
			header := r.Header.Get("Authorization")
			if header == "" {
				logger.Warn("auth: missing authorization header", "path", r.URL.Path)
				return &exception.UnauthorizedError{Msg: "missing authorization header"}
			}

			tokenString, ok := strings.CutPrefix(header, "Bearer ")
			if !ok {
				logger.Warn("auth: invalid authorization header format", "path", r.URL.Path)
				return &exception.UnauthorizedError{Msg: "invalid authorization header"}
			}

			userID, err := jwtService.Parse(tokenString)
			if err != nil {
				logger.Warn("auth: token parse failed", "error", err, "path", r.URL.Path)
				return &exception.UnauthorizedError{Msg: "invalid token"}
			}

			ctx := context.WithValue(r.Context(), auth.UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))

			return nil
		})
	}
}
