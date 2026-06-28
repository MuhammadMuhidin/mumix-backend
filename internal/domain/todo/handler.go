package todo

import (
	"database/sql"
	"encoding/json"
	"mumix-backend/pkg/errors"
	"mumix-backend/pkg/response"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	repo *TodoRepository
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{repo: NewRepository(db)}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/todos", h.ListCreate)
	mux.HandleFunc("/api/todos/", h.Detail)
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

	var completed *bool
	if c := r.URL.Query().Get("completed"); c != "" {
		val := c == "true"
		completed = &val
	}

	todos, total, err := h.repo.FindAll(r.Context(), limit, offset, completed)
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
	response.Success(w, todos, meta)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title    string  `json:"title"`
		Priority string  `json:"priority"`
		DueDate  *string `json:"due_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteJSON(w, errors.BadRequest("Invalid JSON"))
		return
	}

	if req.Title == "" {
		errors.WriteJSON(w, errors.BadRequest("Title is required"))
		return
	}

	if req.Priority == "" {
		req.Priority = "medium"
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		t, err := time.Parse("2006-01-02", *req.DueDate)
		if err == nil {
			dueDate = &t
		}
	}

	todo, err := h.repo.Create(r.Context(), req.Title, req.Priority, dueDate)
	if err != nil {
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}
	response.Created(w, todo)
}

func (h *Handler) Detail(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/todos/")
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
	case http.MethodPatch:
		h.Patch(w, r, id)
	case http.MethodDelete:
		h.Delete(w, r, id)
	default:
		errors.WriteJSON(w, errors.BadRequest("Method not allowed"))
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request, id int) {
	todo, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}
	if todo == nil {
		errors.WriteJSON(w, errors.NotFound("Todo not found"))
		return
	}
	response.Success(w, todo, nil)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request, id int) {
	var req struct {
		Title    string `json:"title"`
		Completed bool  `json:"completed"`
		Priority string `json:"priority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteJSON(w, errors.BadRequest("Invalid JSON"))
		return
	}

	todo, err := h.repo.Update(r.Context(), id, req.Title, req.Completed, req.Priority, nil)
	if err != nil {
		if err == sql.ErrNoRows {
			errors.WriteJSON(w, errors.NotFound("Todo not found"))
			return
		}
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}
	response.Success(w, todo, nil)
}

func (h *Handler) Patch(w http.ResponseWriter, r *http.Request, id int) {
	var fields map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
		errors.WriteJSON(w, errors.BadRequest("Invalid JSON"))
		return
	}

	todo, err := h.repo.Patch(r.Context(), id, fields)
	if err != nil {
		if err == sql.ErrNoRows {
			errors.WriteJSON(w, errors.NotFound("Todo not found"))
			return
		}
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}
	response.Success(w, todo, nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.repo.Delete(r.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			errors.WriteJSON(w, errors.NotFound("Todo not found"))
			return
		}
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}
	response.NoContent(w)
}
