package todo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func Insert(ctx context.Context, db *sql.DB, title string) (Todo, error) {
	var t Todo

	err := db.QueryRowContext(
		ctx,
		`insert into todos (title)
		 values ($1)
		 returning id, title, completed, created_at, updated_at`,
		title,
	).Scan(
		&t.ID,
		&t.Title,
		&t.Completed,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	return t, err
}

func FindAll(ctx context.Context, db *sql.DB) ([]Todo, error) {
	rows, err := db.QueryContext(ctx, `
		select id, title, completed, created_at, updated_at
		from todos
		order by created_at desc`,
	)
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func FindByID(ctx context.Context, db *sql.DB, id string) (Todo, error) {
	var t Todo

	err := db.QueryRowContext(
		ctx,
		`select id, title, completed, created_at, updated_at
		 from todos
		 where id=$1`,
		id,
	).Scan(
		&t.ID,
		&t.Title,
		&t.Completed,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	return t, err
}

func Update(ctx context.Context, db *sql.DB, id string, t Todo) error {
	res, err := db.ExecContext(
		ctx,
		`update todos
		 set title=$1,
		     completed=$2,
		     updated_at=now()
		 where id=$3`,
		t.Title,
		t.Completed,
		id,
	)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func Patch(ctx context.Context, db *sql.DB, id string, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}

	var sets []string
	var args []any
	i := 1

	for k, v := range fields {
		switch k {
		case "title", "completed":
			sets = append(sets, fmt.Sprintf("%s=$%d", k, i))
			args = append(args, v)
			i++
		}
	}

	if len(sets) == 0 {
		return nil
	}

	sets = append(sets, "updated_at=now()")

	query := fmt.Sprintf(
		"update todos set %s where id=$%d",
		strings.Join(sets, ", "),
		i,
	)

	args = append(args, id)

	res, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func Delete(ctx context.Context, db *sql.DB, id string) error {
	res, err := db.ExecContext(
		ctx,
		`delete from todos where id=$1`,
		id,
	)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}