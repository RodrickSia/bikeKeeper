package parkinglot

import "context"

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

type CreateParams struct {
	Name          string
	Address       string
	Type          string
	TotalCapacity int
	OpenTime      string
	CloseTime     string
	ContactPhone  *string
	ManagerName   *string
	Description   *string
}

func (s *Service) Create(ctx context.Context, p CreateParams) (*ParkingLot, error) {
	lot := &ParkingLot{
		Name: p.Name, Address: p.Address, Type: p.Type, Status: "active",
		TotalCapacity: p.TotalCapacity, OpenTime: p.OpenTime, CloseTime: p.CloseTime,
		ContactPhone: p.ContactPhone, ManagerName: p.ManagerName, Description: p.Description,
	}
	if err := s.repo.Create(ctx, lot); err != nil {
		return nil, err
	}
	return lot, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*ParkingLot, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*ParkingLot, error) {
	return s.repo.List(ctx)
}

type UpdateParams struct {
	Name             *string
	Address          *string
	Type             *string
	Status           *string
	TotalCapacity    *int
	CurrentOccupancy *int
	OpenTime         *string
	CloseTime        *string
	ContactPhone     *string
	ManagerName      *string
	Description      *string
}

func (s *Service) Update(ctx context.Context, id string, p UpdateParams) (*ParkingLot, error) {
	lot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p.Name != nil { lot.Name = *p.Name }
	if p.Address != nil { lot.Address = *p.Address }
	if p.Type != nil { lot.Type = *p.Type }
	if p.Status != nil { lot.Status = *p.Status }
	if p.TotalCapacity != nil { lot.TotalCapacity = *p.TotalCapacity }
	if p.CurrentOccupancy != nil { lot.CurrentOccupancy = *p.CurrentOccupancy }
	if p.OpenTime != nil { lot.OpenTime = *p.OpenTime }
	if p.CloseTime != nil { lot.CloseTime = *p.CloseTime }
	if p.ContactPhone != nil { lot.ContactPhone = p.ContactPhone }
	if p.ManagerName != nil { lot.ManagerName = p.ManagerName }
	if p.Description != nil { lot.Description = p.Description }
	if err := s.repo.Update(ctx, lot); err != nil {
		return nil, err
	}
	return lot, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
