package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"mumix-backend/internal/db"
	"mumix-backend/internal/handlers"
	"mumix-backend/internal/middleware"
	"mumix-backend/internal/repositories"
)

func main() {
	r := gin.Default()

	// Health check (public)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	// DB
  dbURL := os.Getenv("SUPABASE_DB_URL")
  if dbURL == "" {
    panic("SUPABASE_DB_URL is required")
  }
  pool := db.NewPostgresPool(dbURL)
  defer pool.Close()

	// Repo & handler
	expenseRepo := repositories.NewExpenseRepo(pool)
	expenseHandler := handlers.NewExpenseHandler(expenseRepo)

	// Protected routes
	auth := r.Group("/api")
	auth.Use(middleware.SupabaseAuth())

	auth.POST("/expenses", expenseHandler.Create)
	auth.GET("/expenses", expenseHandler.GetAll)
	auth.PUT("/expenses/:id", expenseHandler.Update)
	auth.DELETE("/expenses/:id", expenseHandler.Delete)

	// Run server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
	panic(err)
  }
}
