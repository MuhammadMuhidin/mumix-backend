package todo

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
)

func ListCreate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		switch r.Method {
		case http.MethodPost:
			var req struct {
				Title string `json:"title"`
			}
			json.NewDecoder(r.Body).Decode(&req)

			todo, err := Insert(ctx, db, req.Title)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			json.NewEncoder(w).Encode(todo)

		case http.MethodGet:
			todos, err := FindAll(ctx, db)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			json.NewEncoder(w).Encode(todos)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func Detail(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/todos/"):]
		ctx := context.Background()

		switch r.Method {
		case http.MethodGet:
			t, err := FindByID(ctx, db, id)
			if err != nil {
				http.Error(w, "Not found", 404)
				return
			}
			json.NewEncoder(w).Encode(t)

		case http.MethodPut:
			var t Todo
			json.NewDecoder(r.Body).Decode(&t)
			if err := Update(ctx, db, id, t); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		case http.MethodPatch:
			var fields map[string]any
			json.NewDecoder(r.Body).Decode(&fields)
			Patch(ctx, db, id, fields)
			w.WriteHeader(http.StatusNoContent)

		case http.MethodDelete:
			Delete(ctx, db, id)
			w.WriteHeader(http.StatusNoContent)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
