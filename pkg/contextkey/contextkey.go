package contextkey

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

// ContextKey is a type for context keys to avoid collisions
type ContextKey string

const (
	// UserID is the context key for user ID
	UserID ContextKey = "user_id"
	// UserEmail is the context key for user email
	UserEmail ContextKey = "user_email"
	// UserKey is the context key for user claims
	UserKey ContextKey = "user"
	// RequestID is the context key for request ID
	RequestID ContextKey = "request_id"
)

// GenerateID generates a random ID
func GenerateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// GetUserID retrieves user ID from context
func GetUserID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(UserID).(int)
	return id, ok
}
