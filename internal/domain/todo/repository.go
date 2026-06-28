package todo

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Todo struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Completed bool       `json:"completed"`
	Priority  string     `json:"priority"`
	DueDate   *time.Time `json:"due_date,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type TodoRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) FindAll(ctx context.Context, limit, offset int, completed *bool) ([]Todo, int, error) {
	query := "SELECT id, title, completed, priority, due_date, created_at, updated_at FROM todos WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM todos WHERE 1=1"
	args := []interface{}{}

	if completed != nil {
		query += " AND completed = $1"
		countQuery += " AND completed = $1"
		args = append(args, *completed)
	}

	query += " ORDER BY created_at DESC"

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		var dueDate sql.NullTime
		var updatedAt sql.NullTime
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.Priority, &dueDate, &t.CreatedAt, &updatedAt); err != nil {
			return nil, 0, err
		}
		if dueDate.Valid {
			t.DueDate = &dueDate.Time
		}
		if updatedAt.Valid {
			t.UpdatedAt = &updatedAt.Time
		}
		todos = append(todos, t)
	}
	return todos, total, nil
}

func (r *TodoRepository) FindByID(ctx context.Context, id int) (*Todo, error) {
	var t Todo
	var dueDate sql.NullTime
	var updatedAt sql.NullTime
	err := r.db.QueryRowContext(ctx,
		"SELECT id, title, completed, priority, due_date, created_at, updated_at FROM todos WHERE id = $1", id,
	).Scan(&t.ID, &t.Title, &t.Completed, &t.Priority, &dueDate, &t.CreatedAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if dueDate.Valid {
		t.DueDate = &dueDate.Time
	}
	if updatedAt.Valid {
		t.UpdatedAt = &updatedAt.Time
	}
	return &t, nil
}

func (r *TodoRepository) Create(ctx context.Context, title, priority string, dueDate *time.Time) (*Todo, error) {
	var t Todo
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO todos (title, priority, due_date) VALUES ($1, $2, $3) RETURNING id, title, completed, priority, due_date, created_at",
		title, priority, dueDate,
	).Scan(&t.ID, &t.Title, &t.Completed, &t.Priority, &t.DueDate, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TodoRepository) Update(ctx context.Context, id int, title string, completed bool, priority string, dueDate *time.Time) (*Todo, error) {
	res, err := r.db.ExecContext(ctx,
		"UPDATE todos SET title=$1, completed=$2, priority=$3, due_date=$4, updated_at=NOW() WHERE id=$5",
		title, completed, priority, dueDate, id,
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

func (r *TodoRepository) Patch(ctx context.Context, id int, fields map[string]interface{}) (*Todo, error) {
	if len(fields) == 0 {
		return r.FindByID(ctx, id)
	}

	allowed := map[string]bool{"title": true, "completed": true, "priority": true, "due_date": true}
	setClauses := []string{}
	args := []interface{}{}
	i := 1

	for k, v := range fields {
		if allowed[k] {
			setClauses = append(setClauses, fmt.Sprintf("%s=$%d", k, i))
			args = append(args, v)
			i++
		}
	}

	if len(setClauses) == 0 {
		return r.FindByID(ctx, id)
	}

	query := "UPDATE todos SET "
	for j, clause := range setClauses {
		if j > 0 {
			query += ", "
		}
		query += clause
	}
	query += fmt.Sprintf(" WHERE id=$%d", i)
	args = append(args, id)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, sql.ErrNoRows
	}
	return r.FindByID(ctx, id)
}

func (r *TodoRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
