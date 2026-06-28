package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"mumix-backend/internal/auth"
	"mumix-backend/pkg/contextkey"
	"mumix-backend/pkg/errors"
	"net/http"
	"time"
)

// Logger logs HTTP request details
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)

		slog.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"duration", time.Since(start),
			"user_agent", r.UserAgent(),
		)
	})
}

// Recovery recovers from panics and returns a 500 error
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered", "error", rec)
				errors.WriteJSON(w, errors.Internal("Internal Server Error"))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// CORS adds CORS headers to responses
func CORS(origin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequestID adds a request ID to the context
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		w.Header().Set("X-Request-ID", requestID)

		ctx := context.WithValue(r.Context(), contextkey.RequestID, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Authenticate verifies JWT tokens and adds claims to context
func Authenticate(db interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
				tokenString = tokenString[7:]
			}

			claims, err := auth.VerifyToken(tokenString)
			if err != nil {
				errors.WriteJSON(w, errors.Unauthorized("Invalid or expired token"))
				return
			}

			ctx := context.WithValue(r.Context(), contextkey.UserKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves user claims from context
func GetUserFromContext(ctx context.Context) (*auth.UserClaims, bool) {
	claims, ok := ctx.Value(contextkey.UserKey).(*auth.UserClaims)
	return claims, ok
}

func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// statusWriter wraps http.ResponseWriter to capture status code
type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
