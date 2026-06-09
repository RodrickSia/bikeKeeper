package visitor

import (
	"context"
	"fmt"
)

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

type CreateParams struct {
	UserID       string
	VisitorName  string
	VisitorPhone *string
	VehiclePlate string
	ValidDate    string
}

func (s *Service) Create(ctx context.Context, p CreateParams) (*VisitorPass, error) {
	qr := fmt.Sprintf("VISITOR-%s-%s", p.VehiclePlate, p.ValidDate)
	pass := &VisitorPass{
		UserID: p.UserID, VisitorName: p.VisitorName, VisitorPhone: p.VisitorPhone,
		VehiclePlate: p.VehiclePlate, ValidDate: p.ValidDate, Status: "valid", QRCodeData: qr,
	}
	if err := s.repo.Create(ctx, pass); err != nil {
		return nil, err
	}
	return pass, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*VisitorPass, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListByUser(ctx context.Context, userID string) ([]*VisitorPass, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *Service) Cancel(ctx context.Context, id string) error {
	return s.repo.UpdateStatus(ctx, id, "cancelled")
}
