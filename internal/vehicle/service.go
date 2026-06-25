package vehicle

import (
	"context"
	"fmt"
)

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

type CreateParams struct {
	LicensePlate string
	Brand        string
	Model        string
	Color        string
	OwnerID      string
}

func (s *Service) Add(ctx context.Context, p CreateParams) (*Vehicle, error) {
	existing, err := s.repo.FindByPlate(ctx, p.LicensePlate)
	if err != nil {
		return nil, fmt.Errorf("checking existing plate: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("vehicle with plate %s already registered", p.LicensePlate)
	}
	v := &Vehicle{LicensePlate: p.LicensePlate, Brand: p.Brand, Model: p.Model, Color: p.Color, OwnerID: p.OwnerID}
	if err := s.repo.Create(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*Vehicle, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) MyVehicles(ctx context.Context, ownerID string) ([]*Vehicle, error) {
	return s.repo.ListByOwner(ctx, ownerID)
}

func (s *Service) List(ctx context.Context) ([]*Vehicle, error) {
	return s.repo.List(ctx)
}

func (s *Service) FindByPlate(ctx context.Context, plate string) (*Vehicle, error) {
	return s.repo.FindByPlate(ctx, plate)
}

func (s *Service) Remove(ctx context.Context, id, ownerID string) error {
	return s.repo.Deactivate(ctx, id, ownerID)
}
