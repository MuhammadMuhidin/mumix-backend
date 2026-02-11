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

	// Root health check
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("service is up"))
	})
	return mux
}
