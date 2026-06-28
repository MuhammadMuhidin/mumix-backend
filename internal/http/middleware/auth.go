package middleware

import (
	"context"
	"mumix-backend/internal/auth"
	"mumix-backend/pkg/contextkey"
	"mumix-backend/pkg/response"
	"net/http"
	"strings"
)

// Auth creates an authentication middleware
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				resp := response.New(w)
				resp.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header required",
				})
				return
			}

			// Extract Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				resp := response.New(w)
				resp.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization header format",
				})
				return
			}

			// Validate token
			claims, err := auth.VerifyToken(parts[1])
			if err != nil {
				resp := response.New(w)
				resp.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid or expired token",
				})
				return
			}

			// Add claims to context
			ctx := context.WithValue(r.Context(), contextkey.UserID, claims.UserID)
			ctx = context.WithValue(ctx, contextkey.UserEmail, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
