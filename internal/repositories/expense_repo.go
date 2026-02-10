package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Expense struct {
	ID    int64     `json:"id"`
	Name  string    `json:"name"`
	Price *int64    `json:"price,omitempty"`
	Type  *string   `json:"type,omitempty"`
	Date  time.Time `json:"date"`
}

type ExpenseRepo struct {
	db *pgxpool.Pool
}

func NewExpenseRepo(db *pgxpool.Pool) *ExpenseRepo {
	return &ExpenseRepo{db: db}
}

// CREATE
func (r *ExpenseRepo) Create(ctx context.Context, e *Expense) error {
	return r.db.QueryRow(ctx,
		`insert into expense_tracking (name, price, type)
		 values ($1, $2, $3)
		 returning id, date`,
		e.Name, e.Price, e.Type,
	).Scan(&e.ID, &e.Date)
}

// READ ALL
func (r *ExpenseRepo) GetAll(ctx context.Context) ([]Expense, error) {
	rows, err := r.db.Query(ctx,
		`select id, name, price, type, date
		 from expense_tracking
		 order by date desc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Expense
	for rows.Next() {
		var e Expense
		if err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.Price,
			&e.Type,
			&e.Date,
		); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}

// READ BY ID
func (r *ExpenseRepo) GetByID(ctx context.Context, id int64) (*Expense, error) {
	var e Expense
	err := r.db.QueryRow(ctx,
		`select id, name, price, type, date
		 from expense_tracking where id=$1`,
		id,
	).Scan(&e.ID, &e.Name, &e.Price, &e.Type, &e.Date)

	if err != nil {
		return nil, err
	}
	return &e, nil
}

// UPDATE
func (r *ExpenseRepo) Update(ctx context.Context, e *Expense) error {
	_, err := r.db.Exec(ctx,
		`update expense_tracking
		 set name=$1, price=$2, type=$3
		 where id=$4`,
		e.Name, e.Price, e.Type, e.ID,
	)
	return err
}

// DELETE
func (r *ExpenseRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx,
		`delete from expense_tracking where id=$1`,
		id,
	)
	return err
}