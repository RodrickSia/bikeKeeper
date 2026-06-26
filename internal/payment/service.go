package payment

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

func (s *Service) Deposit(ctx context.Context, cardUID string, amount float64) (*Transaction, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("deposit amount must be positive")
	}
	return s.repo.Deposit(ctx, cardUID, amount)
}

func (s *Service) Withdraw(ctx context.Context, cardUID string, amount float64, txType string) (*Transaction, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("withdraw amount must be positive")
	}
	return s.repo.Withdraw(ctx, cardUID, amount, txType)
}

func (s *Service) ChargeParking(ctx context.Context, cardUID string, fee float64, sessionID int64) (*Transaction, error) {
	if fee <= 0 {
		return nil, fmt.Errorf("fee must be positive")
	}
	return s.repo.ChargeParking(ctx, cardUID, fee, sessionID)
}

func (s *Service) ListByCard(ctx context.Context, cardUID string) ([]*Transaction, error) {
	return s.repo.ListByCard(ctx, cardUID)
}

func (s *Service) GetBalance(ctx context.Context, cardUID string) (float64, error) {
	return s.repo.GetBalance(ctx, cardUID)
}
