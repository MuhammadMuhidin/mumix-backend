package handlers

import (
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