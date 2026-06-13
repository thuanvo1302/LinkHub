package middleware

import (
	"context"
	"net/http"
	"strings"

	"linkhub/backend/internal/config"
	"linkhub/backend/internal/pkg/jwt"
	"linkhub/backend/internal/pkg/response"
)

type contextKey string

const userIDContextKey contextKey = "userID"

func Auth(cfg config.Config, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing bearer token", nil)
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwt.Parse(token, cfg.JWTSecret)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token", nil)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) string {
	value, _ := ctx.Value(userIDContextKey).(string)
	return value
}
