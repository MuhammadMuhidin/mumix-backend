package main

import (
	"log"
	"net/http"
	"os"

	"mumix-backend/internal/db"
	httpRouter "mumix-backend/internal/http"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	database, err := db.New(dsn)
	if err != nil {
		log.Fatal(err)
	}

	router := httpRouter.New(database)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
