package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type User struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Phone        string     `json:"phone,omitempty"`
	Role         string     `json:"role"`
	IsActive     bool       `json:"is_active"`
	TokenVersion int        `json:"token_version"`
	Password     string     `json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

type UserRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindAll(ctx context.Context, limit, offset int) ([]User, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE is_active = true").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		"SELECT id, name, email, phone, role, is_active, token_version, created_at, updated_at FROM users WHERE is_active = true ORDER BY id ASC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Role, &u.IsActive, &u.TokenVersion, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int) (*User, error) {
	var u User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, email, phone, role, is_active, token_version, created_at, updated_at FROM users WHERE id = $1", id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Role, &u.IsActive, &u.TokenVersion, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, email, phone, password, role, is_active, token_version, created_at, updated_at FROM users WHERE email = $1", email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Password, &u.Role, &u.IsActive, &u.TokenVersion, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByIDWithPassword(ctx context.Context, id int) (*User, error) {
	var u User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, email, phone, password, role, is_active, token_version, created_at, updated_at FROM users WHERE id = $1", id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Password, &u.Role, &u.IsActive, &u.TokenVersion, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, name, email, phone, password, role string) (*User, error) {
	var u User
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO users (name, email, phone, password, role) VALUES ($1, $2, $3, $4, $5) RETURNING id, name, email, phone, role, is_active, token_version, created_at",
		name, email, phone, password, role,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Role, &u.IsActive, &u.TokenVersion, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	query := "UPDATE users SET "
	args := []interface{}{}
	i := 1
	setClauses := []string{}

	for k, v := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s=$%d", k, i))
		args = append(args, v)
		i++
	}

	if len(setClauses) == 0 {
		return nil
	}

	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(" WHERE id=$%d", i)
	args = append(args, id)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, "UPDATE users SET is_active = false WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepository) HardDelete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
