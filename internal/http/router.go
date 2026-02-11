package http

import (
	"database/sql"
	"net/http"

	"mumix-backend/internal/todo"
)

func New(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/todos", todo.ListCreate(db))
	mux.Handle("/todos/", todo.Detail(db))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("service is up"))
	})

	return mux
}