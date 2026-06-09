package device

import "context"

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) Create(ctx context.Context, d *Device) error {
	return s.repo.Create(ctx, d)
}

func (s *Service) GetByID(ctx context.Context, id string) (*Device, error) {
	d, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	count, _ := s.repo.AlertCount(ctx, id)
	d.AlertCount = count
	return d, nil
}

func (s *Service) List(ctx context.Context) ([]*Device, error) {
	devices, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, d := range devices {
		count, _ := s.repo.AlertCount(ctx, d.ID)
		d.AlertCount = count
	}
	return devices, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id, status string, notes *string) (*Device, error) {
	if err := s.repo.UpdateStatus(ctx, id, status, notes); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *Service) UpdateNotes(ctx context.Context, id, notes string) (*Device, error) {
	if err := s.repo.UpdateNotes(ctx, id, notes); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *Service) ReportFault(ctx context.Context, deviceID, message, severity string) (*Alert, error) {
	a := &Alert{DeviceID: deviceID, Message: message, Severity: severity}
	if err := s.repo.CreateAlert(ctx, a); err != nil {
		return nil, err
	}
	_ = s.repo.UpdateStatus(ctx, deviceID, "warning", nil)
	return a, nil
}

func (s *Service) GetAlerts(ctx context.Context, deviceID *string) ([]*Alert, error) {
	return s.repo.ListAlerts(ctx, deviceID)
}

func (s *Service) ResolveAlert(ctx context.Context, alertID string) error {
	return s.repo.ResolveAlert(ctx, alertID)
}
