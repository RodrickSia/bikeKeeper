package notification

import "context"

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

type CreateParams struct {
	UserID  string
	Title   string
	Message string
	Type    string
}

func (s *Service) Create(ctx context.Context, p CreateParams) (*Notification, error) {
	n := &Notification{UserID: p.UserID, Title: p.Title, Message: p.Message, Type: p.Type}
	if err := s.repo.Create(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) ListByUser(ctx context.Context, userID string) ([]*Notification, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *Service) MarkRead(ctx context.Context, id string) error {
	return s.repo.MarkRead(ctx, id)
}

func (s *Service) MarkAllRead(ctx context.Context, userID string) error {
	return s.repo.MarkAllRead(ctx, userID)
}
