package member

import (
	"context"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

type CreateParams struct {
	StudentID string
	FullName  string
	Phone     *string
}

func (s *Service) Create(ctx context.Context, params CreateParams) (*Member, error) {
	member := &Member{
		StudentID: params.StudentID,
		FullName:  params.FullName,
		Phone:     params.Phone,
	}

	if err := s.repo.Create(ctx, member); err != nil {
		return nil, fmt.Errorf("creating member: %w", err)
	}
	return member, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*Member, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByStudentID(ctx context.Context, studentID string) (*Member, error) {
	return s.repo.GetByStudentID(ctx, studentID)
}

func (s *Service) List(ctx context.Context) ([]*Member, error) {
	return s.repo.List(ctx)
}

type UpdateParams struct {
	ID       string
	FullName *string
	Phone    *string
}

func (s *Service) Update(ctx context.Context, params UpdateParams) (*Member, error) {
	member, err := s.repo.GetByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	if params.FullName != nil {
		member.FullName = *params.FullName
	}
	if params.Phone != nil {
		member.Phone = params.Phone
	}

	if err := s.repo.Update(ctx, member); err != nil {
		return nil, fmt.Errorf("updating member: %w", err)
	}
	return member, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
