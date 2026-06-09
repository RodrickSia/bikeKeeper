package cardrequest

import (
	"context"
	"fmt"
)

type CardRepo interface {
	CreateCard(ctx context.Context, cardUID, cardType string, memberID *string) error
}

type Service struct {
	repo     Repository
	cardRepo CardRepo
}

func NewService(repo Repository, cardRepo CardRepo) *Service {
	return &Service{repo: repo, cardRepo: cardRepo}
}

type CreateParams struct {
	MemberID     string
	VehiclePlate string
	VehicleBrand string
	VehicleModel string
	VehicleColor string
	IDCardNumber string
	Note         *string
}

func (s *Service) Submit(ctx context.Context, p CreateParams) (*CardRequest, error) {
	req := &CardRequest{
		MemberID:     p.MemberID,
		VehiclePlate: p.VehiclePlate,
		VehicleBrand: p.VehicleBrand,
		VehicleModel: p.VehicleModel,
		VehicleColor: p.VehicleColor,
		IDCardNumber: p.IDCardNumber,
		Note:         p.Note,
		Status:       StatusPending,
	}
	if err := s.repo.Create(ctx, req); err != nil {
		return nil, fmt.Errorf("submitting card request: %w", err)
	}
	return req, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*CardRequest, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListByMember(ctx context.Context, memberID string) ([]*CardRequest, error) {
	return s.repo.ListByMember(ctx, memberID)
}

func (s *Service) List(ctx context.Context, status *string) ([]*CardRequest, error) {
	return s.repo.List(ctx, status)
}

func (s *Service) Approve(ctx context.Context, id, cardUID, reviewedBy string) (*CardRequest, error) {
	req, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Status != StatusPending {
		return nil, fmt.Errorf("request is not pending")
	}
	cardType := "monthly"
	memberID := &req.MemberID
	if err := s.cardRepo.CreateCard(ctx, cardUID, cardType, memberID); err != nil {
		return nil, fmt.Errorf("creating card: %w", err)
	}
	uid := cardUID
	rb := reviewedBy
	if err := s.repo.UpdateStatus(ctx, id, StatusApproved, &uid, nil, &rb); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Reject(ctx context.Context, id, reason, reviewedBy string) (*CardRequest, error) {
	req, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Status != StatusPending {
		return nil, fmt.Errorf("request is not pending")
	}
	rb := reviewedBy
	if err := s.repo.UpdateStatus(ctx, id, StatusRejected, nil, &reason, &rb); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Block(ctx context.Context, id string) (*CardRequest, error) {
	if err := s.repo.UpdateStatus(ctx, id, StatusBlocked, nil, nil, nil); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}
