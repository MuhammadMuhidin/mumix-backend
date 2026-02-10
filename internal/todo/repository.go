package todo

import (
	"context"
	"database/sql"
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
	rows, err := db.QueryContext(
		ctx,
		`select id, title, completed, created_at, updated_at
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
	_, err := db.ExecContext(
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

	return err
}

func Patch(ctx context.Context, db *sql.DB, id string, fields map[string]any) error {
	if v, ok := fields["title"]; ok {
		_, _ = db.ExecContext(
			ctx,
			`update todos
			 set title=$1,
			     updated_at=now()
			 where id=$2`,
			v,
			id,
		)
	}

	if v, ok := fields["completed"]; ok {
		_, _ = db.ExecContext(
			ctx,
			`update todos
			 set completed=$1,
			     updated_at=now()
			 where id=$2`,
			v,
			id,
		)
	}

	return nil
}

func Delete(ctx context.Context, db *sql.DB, id string) error {
	_, err := db.ExecContext(
		ctx,
		`delete from todos
		 where id=$1`,
		id,
	)

	return err
}
