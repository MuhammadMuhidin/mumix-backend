package expense

import (
	"database/sql"
	"encoding/json"
	"mumix-backend/pkg/errors"
	"mumix-backend/pkg/response"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	repo *ExpenseRepository
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{repo: NewRepository(db)}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/expenses", h.ListCreate)
	mux.HandleFunc("/api/expenses/", h.Detail)
	mux.HandleFunc("/api/expenses/totals", h.Totals)
	mux.HandleFunc("/api/expenses/export", h.ExportCSV)
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
	expenses, err := h.repo.FindAll(r.Context())
	if err != nil {
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}

	var cumSisa, cumTabungan int
	var responses []ExpenseResponse
	for _, e := range expenses {
		resp := ToExpenseResponse(&e, cumSisa, cumTabungan)
		cumSisa = parseRupiah(resp.TotalSisa)
		cumTabungan = parseRupiah(resp.TotalTabungan)
		responses = append(responses, resp)
	}

	response.Success(w, responses, nil)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input ExpenseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.WriteJSON(w, errors.BadRequest("Invalid JSON"))
		return
	}

	if input.Month == "" {
		errors.WriteJSON(w, errors.BadRequest("Month is required"))
		return
	}

	e := input.ToExpense()
	created, err := h.repo.Create(r.Context(), e)
	if err != nil {
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}

	var cumSisa, cumTabungan int
	resp := ToExpenseResponse(created, cumSisa, cumTabungan)
	response.Created(w, resp)
}

func (h *Handler) Detail(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/expenses/")
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
	e, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}
	if e == nil {
		errors.WriteJSON(w, errors.NotFound("Expense not found"))
		return
	}

	var cumSisa, cumTabungan int
	resp := ToExpenseResponse(e, cumSisa, cumTabungan)
	response.Success(w, resp, nil)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request, id int) {
	var input ExpenseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.WriteJSON(w, errors.BadRequest("Invalid JSON"))
		return
	}

	e := input.ToExpense()
	updated, err := h.repo.Update(r.Context(), id, e)
	if err != nil {
		if err == sql.ErrNoRows {
			errors.WriteJSON(w, errors.NotFound("Expense not found"))
			return
		}
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}

	var cumSisa, cumTabungan int
	resp := ToExpenseResponse(updated, cumSisa, cumTabungan)
	response.Success(w, resp, nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.repo.Delete(r.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			errors.WriteJSON(w, errors.NotFound("Expense not found"))
			return
		}
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}
	response.NoContent(w)
}

func (h *Handler) Totals(w http.ResponseWriter, r *http.Request) {
	totals, err := h.repo.GetTotals(r.Context())
	if err != nil {
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}

	formatted := map[string]string{}
	for k, v := range totals {
		formatted[k] = formatRupiah(v)
	}

	response.Success(w, formatted, nil)
}

func (h *Handler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	expenses, err := h.repo.FindAll(r.Context())
	if err != nil {
		errors.WriteJSON(w, errors.Internal(err.Error()))
		return
	}

	var cumSisa, cumTabungan int
	var lines []string
	lines = append(lines, "Tahun,Bulan,Pendapatan Utama,Pendapatan Lainnya,Transfer ke Nda,Transfer ke Mimi,Kartu Kredit,Tabungan Tetap,Pengeluaran Lainnya,Sisa Bulan Ini,Total Sisa,Total Tabungan")
	for _, e := range expenses {
		resp := ToExpenseResponse(&e, cumSisa, cumTabungan)
		cumSisa = parseRupiah(resp.TotalSisa)
		cumTabungan = parseRupiah(resp.TotalTabungan)

		line := strings.Join([]string{
			e.Year, e.Month,
			resp.IncomeMain, resp.IncomeOther,
			resp.TransferToNda, resp.TransferToMimi,
			resp.CreditCard, resp.FixedSavings,
			resp.OtherExpense,
			resp.SisaBulanIni, resp.TotalSisa, resp.TotalTabungan,
		}, ",")
		lines = append(lines, line)
	}

	csv := strings.Join(lines, "\n")

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="expenses.csv"`)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(csv))
}
