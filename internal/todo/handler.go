package todo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func ListCreate(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodPost:
			var req struct {
				Title string `json:"title"`
			}

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}

			if req.Title == "" {
				http.Error(w, "title is required", http.StatusBadRequest)
				return
			}

			todo, err := Insert(r.Context(), db, req.Title)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(todo)

		case http.MethodGet:
			todos, err := FindAll(r.Context(), db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			_ = json.NewEncoder(w).Encode(todos)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

func Detail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := strings.TrimPrefix(r.URL.Path, "/todos/")
		if id == "" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			todo, err := FindByID(r.Context(), db, id)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					http.Error(w, "not found", http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			_ = json.NewEncoder(w).Encode(todo)

		case http.MethodPut:
			var t Todo
			if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}

			if err := Update(r.Context(), db, id, t); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					http.Error(w, "not found", http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		case http.MethodPatch:
			var fields map[string]any
			if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}

			if err := Patch(r.Context(), db, id, fields); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					http.Error(w, "not found", http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		case http.MethodDelete:
			if err := Delete(r.Context(), db, id); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					http.Error(w, "not found", http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}