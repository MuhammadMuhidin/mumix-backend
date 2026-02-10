package main

import (
	"log"
	"net/http"
	"os"

	"mumix-backend/internal/db"
	httpRouter "mumix-backend/internal/http"
)

func main() {
	dsn := os.Getenv("SPBASE_DB_URL")
	if dsn == "" {
		log.Fatal("SPBASE is required")
	}

	database, err := db.New(dsn)
	if err != nil {
		log.Fatal(err)
	}

	router := httpRouter.New(database)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
	log.Println(len(os.Getenv("SPBASE_DB_URL")))
	log.Println("DB URI:", os.Getenv("SPBASE_DB_URL"))
}
