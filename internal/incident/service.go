package incident

import "context"

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

type CreateParams struct {
	ReportedBy   string
	ReporterName string
	VehiclePlate *string
	Type         string
	Description  string
	Location     *string
}

func (s *Service) Report(ctx context.Context, p CreateParams) (*Incident, error) {
	inc := &Incident{
		ReportedBy: p.ReportedBy, ReporterName: p.ReporterName,
		VehiclePlate: p.VehiclePlate, Type: p.Type,
		Description: p.Description, Location: p.Location, Status: "open",
	}
	if err := s.repo.Create(ctx, inc); err != nil {
		return nil, err
	}
	return inc, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*Incident, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, staffID *string) ([]*Incident, error) {
	return s.repo.List(ctx, staffID)
}

func (s *Service) Resolve(ctx context.Context, id, note string) (*Incident, error) {
	if err := s.repo.Resolve(ctx, id, note); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Escalate(ctx context.Context, id string) (*Incident, error) {
	if err := s.repo.Escalate(ctx, id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}
