package parkingsession

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

type CheckInParams struct {
	CardUID         string
	PlateIn         *string
	ImgPlateInPath  *string
	ImgPersonInPath *string
}

type CheckOutParams struct {
	PlateOut         *string
	ImgPlateOutPath  *string
	ImgPersonOutPath *string
	Cost             decimal.Decimal
	IsWarning        bool
}

func (s *Service) CheckIn(ctx context.Context, params CheckInParams) (*ParkingSession, error) {
	ongoing, err := s.repo.GetOngoingSessionByCard(ctx, params.CardUID)
	if err != nil {
		return nil, fmt.Errorf("checking ongoing session: %w", err)
	}
	if ongoing != nil {
		return nil, fmt.Errorf("card %s already has an ongoing session (id: %d)", params.CardUID, ongoing.ID)
	}

	session := &ParkingSession{
		CardUID:         params.CardUID,
		PlateIn:         params.PlateIn,
		ImgPlateInPath:  params.ImgPlateInPath,
		ImgPersonInPath: params.ImgPersonInPath,
		Status:          "ongoing",
	}

	if err := s.repo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("creating parking session: %w", err)
	}

	return session, nil
}

func (s *Service) CheckOut(ctx context.Context, id int64, params CheckOutParams) error {
	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("fetching session: %w", err)
	}
	if session.Status != "ongoing" {
		return fmt.Errorf("session %d is already completed", id)
	}
	if err := s.repo.CheckOut(ctx, id, params); err != nil {
		return fmt.Errorf("checking out session: %w", err)
	}
	return nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*ParkingSession, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListByCard(ctx context.Context, cardUID string) ([]*ParkingSession, error) {
	return s.repo.ListByCard(ctx, cardUID)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
