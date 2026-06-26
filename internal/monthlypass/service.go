package monthlypass

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, p CreateParams) (*MonthlyPass, error) {
	return s.repo.Create(ctx, p)
}

func (s *Service) ListByUser(ctx context.Context, userID string) ([]*MonthlyPass, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *Service) ToggleAutoRenew(ctx context.Context, id string) (*MonthlyPass, error) {
	return s.repo.ToggleAutoRenew(ctx, id)
}
