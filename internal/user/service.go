package user

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

type CreateParams struct {
	Email    string
	Password string
	Role     string
	MemberID *string
}

func (s *Service) Create(ctx context.Context, params CreateParams) (*User, error) {
	switch params.Role {
	case RoleStudent, RoleStaff, RoleFaculty, RoleAdmin:
	default:
		return nil, fmt.Errorf("invalid role: %s", params.Role)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := &User{
		Email:        params.Email,
		PasswordHash: string(hash),
		Role:         params.Role,
		MemberID:     params.MemberID,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *Service) List(ctx context.Context) ([]*User, error) {
	return s.repo.List(ctx)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
