package auth

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"mumix-backend/pkg/contextkey"
	"mumix-backend/pkg/errors"
	"net/http"
	"os"
)

type Handler struct {
	db *sql.DB
}

func NewAuthHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/register", h.Register)
	mux.HandleFunc("/api/auth/login", h.Login)
	mux.HandleFunc("/api/auth/logout", h.Logout)
	mux.HandleFunc("/api/auth/me", h.Me)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.WriteJSON(w, errors.BadRequest("Method not allowed"))
		return
	}

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteJSON(w, errors.BadRequest("Invalid JSON"))
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		errors.WriteJSON(w, errors.BadRequest("Name, email, and password are required"))
		return
	}

	var exists int
	err := h.db.QueryRowContext(r.Context(),
		"SELECT COUNT(*) FROM users WHERE email = $1", req.Email,
	).Scan(&exists)
	if err != nil {
		errors.WriteJSON(w, errors.Internal("Database error"))
		return
	}
	if exists > 0 {
		errors.WriteJSON(w, errors.Conflict("Email already exists"))
		return
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		errors.WriteJSON(w, errors.Internal("Failed to hash password"))
		return
	}

	var userID int
	var userName, userEmail, userRole string
	err = h.db.QueryRowContext(r.Context(),
		"INSERT INTO users (name, email, phone, password, role) VALUES ($1, $2, $3, $4, 'user') RETURNING id, name, email, role",
		req.Name, req.Email, req.Phone, hashedPassword,
	).Scan(&userID, &userName, &userEmail, &userRole)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		errors.WriteJSON(w, errors.Internal("Failed to create user"))
		return
	}

	token, err := GenerateToken(userID, userRole, 0)
	if err != nil {
		slog.Error("failed to generate token", "error", err)
		errors.WriteJSON(w, errors.Internal("Failed to generate token"))
		return
	}

	setTokenCookie(w, token, os.Getenv("NODE_ENV") == "production")

	resp := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"user": map[string]interface{}{
				"id":    userID,
				"name":  userName,
				"email": userEmail,
				"role":  userRole,
			},
			"token": token,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.WriteJSON(w, errors.BadRequest("Method not allowed"))
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteJSON(w, errors.BadRequest("Invalid JSON"))
		return
	}

	var userID int
	var userName, userEmail, userRole, hashedPassword string
	var isActive bool
	var tokenVersion int
	err := h.db.QueryRowContext(r.Context(),
		"SELECT id, name, email, password, role, is_active, token_version FROM users WHERE email = $1",
		req.Email,
	).Scan(&userID, &userName, &userEmail, &hashedPassword, &userRole, &isActive, &tokenVersion)
	if err != nil {
		slog.Warn("login attempt for non-existent user", "email", req.Email)
		errors.WriteJSON(w, errors.Unauthorized("Invalid email or password"))
		return
	}

	if !isActive {
		errors.WriteJSON(w, errors.Unauthorized("Account is disabled"))
		return
	}

	if !CheckPassword(req.Password, hashedPassword) {
		slog.Warn("login failed - wrong password", "email", req.Email, "ip", r.RemoteAddr)
		errors.WriteJSON(w, errors.Unauthorized("Invalid email or password"))
		return
	}

	token, err := GenerateToken(userID, userRole, tokenVersion)
	if err != nil {
		slog.Error("failed to generate token", "error", err)
		errors.WriteJSON(w, errors.Internal("Failed to generate token"))
		return
	}

	setTokenCookie(w, token, os.Getenv("NODE_ENV") == "production")

	resp := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"user": map[string]interface{}{
				"id":    userID,
				"name":  userName,
				"email": userEmail,
				"role":  userRole,
			},
			"token": token,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(contextkey.UserKey).(*UserClaims)
	if !ok {
		errors.WriteJSON(w, errors.Unauthorized("Authentication required"))
		return
	}

	var userName, userEmail, userRole string
	var isActive bool
	err := h.db.QueryRowContext(r.Context(),
		"SELECT name, email, role, is_active FROM users WHERE id = $1",
		claims.UserID,
	).Scan(&userName, &userEmail, &userRole, &isActive)
	if err != nil {
		errors.WriteJSON(w, errors.NotFound("User not found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":        claims.UserID,
			"name":      userName,
			"email":     userEmail,
			"role":      userRole,
			"is_active": isActive,
		},
	})
}

func setTokenCookie(w http.ResponseWriter, token string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   secure,
		MaxAge:   86400,
	})
}
