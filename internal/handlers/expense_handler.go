package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"mumix-backend/internal/repositories"
)

type ExpenseHandler struct {
	repo *repositories.ExpenseRepo
}

func NewExpenseHandler(r *repositories.ExpenseRepo) *ExpenseHandler {
	return &ExpenseHandler{repo: r}
}

// POST /expenses
func (h *ExpenseHandler) Create(c *gin.Context) {
	var req struct {
		Name  string  `json:"name" binding:"required"`
		Price *int64  `json:"price"`
		Type  *string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	exp := &repositories.Expense{
		Name:  req.Name,
		Price: req.Price,
		Type:  req.Type,
	}

	if err := h.repo.Create(c.Request.Context(), exp); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, exp)
}

// GET /expenses
func (h *ExpenseHandler) GetAll(c *gin.Context) {
	list, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

// PUT /expenses/:id
func (h *ExpenseHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	var req repositories.Expense
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	req.ID = id

	if err := h.repo.Update(c.Request.Context(), &req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, req)
}

// DELETE /expenses/:id
func (h *ExpenseHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(204)
}
