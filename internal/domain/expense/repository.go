package expense

import (
	"context"
	"database/sql"
	"time"
)

type Expense struct {
	ID             int        `json:"id"`
	Year           string     `json:"year"`
	Month          string     `json:"month"`
	IncomeMain     int        `json:"income_main"`
	IncomeOther    int        `json:"income_other"`
	TransferToNda  int        `json:"transfer_to_nda"`
	TransferToMimi int        `json:"transfer_to_mimi"`
	CreditCard     int        `json:"credit_card"`
	FixedSavings   int        `json:"fixed_savings"`
	OtherExpense   int        `json:"other_expense"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

type ExpenseRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) FindAll(ctx context.Context) ([]Expense, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, year, month, income_main, income_other, transfer_to_nda, transfer_to_mimi, credit_card, fixed_savings, other_expense, created_at, updated_at FROM expenses ORDER BY year DESC, month_order(month) DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var e Expense
		var updatedAt sql.NullTime
		if err := rows.Scan(&e.ID, &e.Year, &e.Month, &e.IncomeMain, &e.IncomeOther, &e.TransferToNda, &e.TransferToMimi, &e.CreditCard, &e.FixedSavings, &e.OtherExpense, &e.CreatedAt, &updatedAt); err != nil {
			return nil, err
		}
		if updatedAt.Valid {
			e.UpdatedAt = &updatedAt.Time
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

func (r *ExpenseRepository) FindByID(ctx context.Context, id int) (*Expense, error) {
	var e Expense
	var updatedAt sql.NullTime
	err := r.db.QueryRowContext(ctx,
		"SELECT id, year, month, income_main, income_other, transfer_to_nda, transfer_to_mimi, credit_card, fixed_savings, other_expense, created_at, updated_at FROM expenses WHERE id = $1", id,
	).Scan(&e.ID, &e.Year, &e.Month, &e.IncomeMain, &e.IncomeOther, &e.TransferToNda, &e.TransferToMimi, &e.CreditCard, &e.FixedSavings, &e.OtherExpense, &e.CreatedAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if updatedAt.Valid {
		e.UpdatedAt = &updatedAt.Time
	}
	return &e, nil
}

func (r *ExpenseRepository) Create(ctx context.Context, e *Expense) (*Expense, error) {
	var id int
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO expenses (year, month, income_main, income_other, transfer_to_nda, transfer_to_mimi, credit_card, fixed_savings, other_expense) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
		e.Year, e.Month, e.IncomeMain, e.IncomeOther, e.TransferToNda, e.TransferToMimi, e.CreditCard, e.FixedSavings, e.OtherExpense,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *ExpenseRepository) Update(ctx context.Context, id int, e *Expense) (*Expense, error) {
	res, err := r.db.ExecContext(ctx,
		"UPDATE expenses SET year=$1, month=$2, income_main=$3, income_other=$4, transfer_to_nda=$5, transfer_to_mimi=$6, credit_card=$7, fixed_savings=$8, other_expense=$9, updated_at=NOW() WHERE id=$10",
		e.Year, e.Month, e.IncomeMain, e.IncomeOther, e.TransferToNda, e.TransferToMimi, e.CreditCard, e.FixedSavings, e.OtherExpense, id,
	)
	if err != nil {
		return nil, err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, sql.ErrNoRows
	}
	return r.FindByID(ctx, id)
}

func (r *ExpenseRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM expenses WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ExpenseRepository) GetTotals(ctx context.Context) (map[string]int, error) {
	var totals struct {
		IncomeMain    int
		IncomeOther   int
		TransferToNda int
		TransferToMimi int
		CreditCard    int
		FixedSavings  int
		OtherExpense  int
	}

	err := r.db.QueryRowContext(ctx,
		"SELECT COALESCE(SUM(income_main), 0), COALESCE(SUM(income_other), 0), COALESCE(SUM(transfer_to_nda), 0), COALESCE(SUM(transfer_to_mimi), 0), COALESCE(SUM(credit_card), 0), COALESCE(SUM(fixed_savings), 0), COALESCE(SUM(other_expense), 0) FROM expenses",
	).Scan(&totals.IncomeMain, &totals.IncomeOther, &totals.TransferToNda, &totals.TransferToMimi, &totals.CreditCard, &totals.FixedSavings, &totals.OtherExpense)

	return map[string]int{
		"pendapatan_utama":    totals.IncomeMain,
		"pendapatan_lainnya":  totals.IncomeOther,
		"transfer_ke_nda":    totals.TransferToNda,
		"transfer_ke_mimi":   totals.TransferToMimi,
		"kartu_kredit":       totals.CreditCard,
		"tabungan_tetap":     totals.FixedSavings,
		"pengeluaran_lainnya": totals.OtherExpense,
	}, err
}
