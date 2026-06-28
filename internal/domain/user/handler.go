package user

import (
	"database/sql"
	"encoding/json"
	"mumix-backend/pkg/errors"
	"mumix-backend/internal/http/middleware"
	"mumix-backend/pkg/response"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{service: NewService(db)}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/users", h.ListCreate)
	mux.HandleFunc("/api/users/", h.Detail)
	mux.HandleFunc("/api/users/change-password", h.ChangePassword)
}

func (h *Handler) ListCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.List(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		errors.WriteJSON(w, errors.BadRequest("Method not allowed"))
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	users, total, err := h.service.GetAll(r.Context(), limit, offset)
	if err != nil {
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}

	meta := map[string]interface{}{
		"page":       page,
		"limit":      limit,
		"total":      total,
		"total_page": (total + limit - 1) / limit,
	}
	response.Success(w, users, meta)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.WriteJSON(w, errors.BadRequest("Invalid JSON"))
		return
	}

	if input.Name == "" || input.Email == "" || input.Password == "" {
		errors.WriteJSON(w, errors.BadRequest("Name, email, and password are required"))
		return
	}

	user, err := h.service.Create(r.Context(), input)
	if err != nil {
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}

	response.Created(w, user)
}

func (h *Handler) Detail(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/api/users/"):]
	id, err := strconv.Atoi(path)
	if err != nil {
		errors.WriteJSON(w, errors.BadRequest("Invalid ID"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r, id)
	case http.MethodPut:
		h.Update(w, r, id)
	case http.MethodDelete:
		h.Delete(w, r, id)
	default:
		errors.WriteJSON(w, errors.BadRequest("Method not allowed"))
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request, id int) {
	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}
	response.Success(w, user, nil)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request, id int) {
	var input UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.WriteJSON(w, errors.BadRequest("Invalid JSON"))
		return
	}

	if err := h.service.Update(r.Context(), id, input); err != nil {
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}

	response.NoContent(w)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, id int) {
	hard := r.URL.Query().Get("hard") == "true"
	if err := h.service.Delete(r.Context(), id, hard); err != nil {
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}
	response.NoContent(w)
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.WriteJSON(w, errors.BadRequest("Method not allowed"))
		return
	}

	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		errors.WriteJSON(w, errors.Unauthorized("Authentication required"))
		return
	}

	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.WriteJSON(w, errors.BadRequest("Invalid JSON"))
		return
	}

	if err := h.service.ChangePassword(r.Context(), claims.UserID, input.OldPassword, input.NewPassword); err != nil {
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}

	response.Message(w, "Password changed successfully", http.StatusOK)
}
