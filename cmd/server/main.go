package main

import (
	"log"
	"net/http"
	"os"

	"mumix-backend/internal/db"
	httpRouter "mumix-backend/internal/http"
)

func main() {
	dsn := os.Getenv("KOYEBDB_URI")
	if dsn == "" {
		log.Fatal("KOYEBDB_URI is required")
	}

	database, err := db.New(dsn)
	if err != nil {
		log.Fatalf("db init error: %v", err)
	}
	defer database.Close()

	router := httpRouter.New(database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}