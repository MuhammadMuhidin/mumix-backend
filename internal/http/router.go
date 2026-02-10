package http

import (
	"database/sql"
	"net/http"

	"mumix-backend/internal/todo"
)

func New(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/todos", todo.ListCreate(db))
	mux.HandleFunc("/todos/", todo.Detail(db))
	return mux
}
