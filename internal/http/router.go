package http

import (
	"database/sql"
	"net/http"

	"mumix-backend/internal/auth"
	"mumix-backend/internal/domain/expense"
	"mumix-backend/internal/domain/todo"
	"mumix-backend/internal/domain/user"
	"mumix-backend/internal/http/middleware"
)

// New creates the main HTTP handler
func New(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"down"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Auth routes (public)
	authH := auth.NewAuthHandler(db)
	authH.RegisterRoutes(mux)

	// Protected API routes
	apiMux := http.NewServeMux()

	todoH := todo.NewHandler(db)
	todoH.RegisterRoutes(apiMux)

	expenseH := expense.NewHandler(db)
	expenseH.RegisterRoutes(apiMux)

	userH := user.NewHandler(db)
	userH.RegisterRoutes(apiMux)

	protected := middleware.Authenticate(db)(apiMux)
	mux.Handle("/api/", protected)

	// Apply global middleware
	var handler http.Handler = mux
	handler = middleware.Recovery(handler)
	handler = middleware.CORS("http://localhost:3000")(handler)
	handler = middleware.Logger(handler)
	handler = middleware.RequestID(handler)

	return handler
}
