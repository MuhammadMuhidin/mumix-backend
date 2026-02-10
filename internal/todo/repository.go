package todo

import (
	"context"
	"database/sql"
)

func Insert(ctx context.Context, db *sql.DB, title string) (Todo, error) {
	var t Todo
	err := db.QueryRowContext(ctx,
		`insert into todos (title)
		 values ($1)
		 returning id, title, completed, created_at, updated_at`,
		title,
	).Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func FindAll(ctx context.Context, db *sql.DB) ([]Todo, error) {
	rows, err := db.QueryContext(ctx,
		`select id, title, completed, created_at, updated_at
		 from todos
		 order by created_at desc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Completed,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}
