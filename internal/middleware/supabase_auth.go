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

		client.SetAuth(token)
		user, err := client.GetUser()
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
