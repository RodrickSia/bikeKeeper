package shift

import (
	"context"
	"fmt"
)

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

type CreateParams struct {
	Name      string
	Type      string
	StartTime string
	EndTime   string
	Date      string
	Notes     *string
}

func (s *Service) Create(ctx context.Context, p CreateParams) (*Shift, error) {
	sh := &Shift{Name: p.Name, Type: p.Type, StartTime: p.StartTime, EndTime: p.EndTime, Date: p.Date, Notes: p.Notes, Status: "scheduled"}
	if err := s.repo.Create(ctx, sh); err != nil {
		return nil, err
	}
	sh.StaffIDs = []string{}
	sh.StaffNames = []string{}
	return sh, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*Shift, error) {
	sh, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	ids, names, err := s.repo.GetAssignments(ctx, id)
	if err != nil {
		return nil, err
	}
	sh.StaffIDs = ids
	sh.StaffNames = names
	return sh, nil
}

func (s *Service) List(ctx context.Context, from, to, staffID *string) ([]*Shift, error) {
	shifts, err := s.repo.List(ctx, from, to, staffID)
	if err != nil {
		return nil, err
	}
	for _, sh := range shifts {
		ids, names, _ := s.repo.GetAssignments(ctx, sh.ID)
		sh.StaffIDs = ids
		sh.StaffNames = names
	}
	return shifts, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id, status string) (*Shift, error) {
	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *Service) UpdateNotes(ctx context.Context, id string, notes *string) (*Shift, error) {
	if err := s.repo.UpdateNotes(ctx, id, notes); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) AssignStaff(ctx context.Context, shiftID, userID string, assignedBy *string) (*Shift, error) {
	ids, _, _ := s.repo.GetAssignments(ctx, shiftID)
	for _, id := range ids {
		if id == userID {
			return nil, fmt.Errorf("staff already assigned to this shift")
		}
	}
	if err := s.repo.AssignStaff(ctx, shiftID, userID, assignedBy); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, shiftID)
}

func (s *Service) RemoveStaff(ctx context.Context, shiftID, userID string) (*Shift, error) {
	if err := s.repo.RemoveStaff(ctx, shiftID, userID); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, shiftID)
}
