package support

import "context"

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

type CreateParams struct {
	UserID      string
	UserName    string
	Category    string
	Subject     string
	Description string
}

func (s *Service) Create(ctx context.Context, p CreateParams) (*Ticket, error) {
	t := &Ticket{UserID: p.UserID, UserName: p.UserName, Category: p.Category, Subject: p.Subject, Description: p.Description, Status: "open"}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	t.Responses = []Response{}
	return t, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*Ticket, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resps, err := s.repo.GetResponses(ctx, id)
	if err != nil {
		return nil, err
	}
	t.Responses = resps
	return t, nil
}

func (s *Service) List(ctx context.Context, userID *string) ([]*Ticket, error) {
	return s.repo.List(ctx, userID)
}

func (s *Service) UpdateStatus(ctx context.Context, id, status string) (*Ticket, error) {
	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *Service) AddResponse(ctx context.Context, ticketID, senderID, senderName, message string, isAdmin bool) (*Ticket, error) {
	resp := &Response{TicketID: ticketID, SenderID: senderID, SenderName: senderName, Message: message, IsAdmin: isAdmin}
	if err := s.repo.AddResponse(ctx, resp); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, ticketID)
}
