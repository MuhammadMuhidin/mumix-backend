package user

import (
	"context"
	"database/sql"
	"mumix-backend/internal/auth"
	"mumix-backend/pkg/errors"
)

type Service struct {
	repo *UserRepository
}

func NewService(db *sql.DB) *Service {
	return &Service{repo: NewRepository(db)}
}

type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UpdateUserInput struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
	Role     *string `json:"role,omitempty"`
}

func (s *Service) Create(ctx context.Context, input CreateUserInput) (*User, error) {
	existing, _ := s.repo.FindByEmail(ctx, input.Email)
	if existing != nil {
		return nil, errors.Conflict("Email already exists")
	}

	if input.Role == "" {
		input.Role = "user"
	}

	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, errors.Internal("Failed to hash password")
	}

	return s.repo.Create(ctx, input.Name, input.Email, input.Phone, hashedPassword, input.Role)
}

func (s *Service) GetByID(ctx context.Context, id int) (*User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NotFound("User not found")
	}
	return user, nil
}

func (s *Service) GetAll(ctx context.Context, limit, offset int) ([]User, int, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *Service) Update(ctx context.Context, id int, input UpdateUserInput) error {
	updates := map[string]interface{}{}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.Phone != nil {
		updates["phone"] = *input.Phone
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
		user, _ := s.repo.FindByID(ctx, id)
		if user != nil {
			updates["token_version"] = user.TokenVersion + 1
		}
	}
	if input.Role != nil {
		updates["role"] = *input.Role
	}

	if len(updates) == 0 {
		return nil
	}

	return s.repo.Update(ctx, id, updates)
}

func (s *Service) Delete(ctx context.Context, id int, hard bool) error {
	if hard {
		return s.repo.HardDelete(ctx, id)
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) ChangePassword(ctx context.Context, id int, oldPassword, newPassword string) error {
	user, err := s.repo.FindByIDWithPassword(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.NotFound("User not found")
	}

	if !auth.CheckPassword(oldPassword, user.Password) {
		return errors.Unauthorized("Current password is incorrect")
	}

	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return errors.Internal("Failed to hash password")
	}

	return s.repo.Update(ctx, id, map[string]interface{}{
		"password":      hashedPassword,
		"token_version": user.TokenVersion + 1,
	})
}
