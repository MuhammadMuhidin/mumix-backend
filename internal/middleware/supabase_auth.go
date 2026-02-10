package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/gotrue-go"
)

func SupabaseAuth() gin.HandlerFunc {
	client := gotrue.New(
		os.Getenv("SUPABASE_URL"),
		os.Getenv("SUPABASE_ROLE_KEY"),
	)

	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		user, err := client.GetUser(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		// simpan user ke context
		c.Set("user", user)
		c.Next()
	}
}