package main

import (
	"log"
	"net/http"
	"os"

	"mumix-backend/internal/db"
	httpRouter "mumix-backend/internal/http"
)

func main() {
	dsn := os.Getenv("SPBASE_URI")
	if dsn == "" {
		log.Fatal("SPBASE_URI is required")
	}

	database, err := db.New(dsn)
	if err != nil {
		log.Fatal(err)
	}

	router := httpRouter.New(database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback local
		}

		log.Println("Server running on :" + port)
		log.Fatal(http.ListenAndServe(":"+port, router))
}
