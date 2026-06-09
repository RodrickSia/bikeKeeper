package device

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, d *Device) error
	GetByID(ctx context.Context, id string) (*Device, error)
	List(ctx context.Context) ([]*Device, error)
	UpdateStatus(ctx context.Context, id, status string, notes *string) error
	UpdateNotes(ctx context.Context, id string, notes string) error
	CreateAlert(ctx context.Context, a *Alert) error
	ListAlerts(ctx context.Context, deviceID *string) ([]*Alert, error)
	ResolveAlert(ctx context.Context, alertID string) error
	AlertCount(ctx context.Context, deviceID string) (int, error)
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, d *Device) error {
	const q = `INSERT INTO devices (name,type,location_label,ip_address,firmware_version,notes,installed_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, last_seen`
	return r.db.QueryRowContext(ctx, q, d.Name, d.Type, d.LocationLabel, d.IPAddress, d.FirmwareVersion, d.Notes, d.InstalledAt).
		Scan(&d.ID, &d.LastSeen)
}

func (r *repository) GetByID(ctx context.Context, id string) (*Device, error) {
	d := &Device{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id,name,type,location_label,status,ip_address,firmware_version,notes,installed_at,last_seen FROM devices WHERE id=$1`, id).
		Scan(&d.ID, &d.Name, &d.Type, &d.LocationLabel, &d.Status, &d.IPAddress, &d.FirmwareVersion, &d.Notes, &d.InstalledAt, &d.LastSeen)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("device not found")
	}
	return d, err
}

func (r *repository) List(ctx context.Context) ([]*Device, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id,name,type,location_label,status,ip_address,firmware_version,notes,installed_at,last_seen FROM devices ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var devices []*Device
	for rows.Next() {
		d := &Device{}
		if err := rows.Scan(&d.ID, &d.Name, &d.Type, &d.LocationLabel, &d.Status, &d.IPAddress, &d.FirmwareVersion, &d.Notes, &d.InstalledAt, &d.LastSeen); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, rows.Err()
}

func (r *repository) UpdateStatus(ctx context.Context, id, status string, notes *string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE devices SET status=$2, notes=COALESCE($3,notes), last_seen=NOW() WHERE id=$1`, id, status, notes)
	return err
}

func (r *repository) UpdateNotes(ctx context.Context, id string, notes string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE devices SET notes=$2 WHERE id=$1`, id, notes)
	return err
}

func (r *repository) CreateAlert(ctx context.Context, a *Alert) error {
	const q = `INSERT INTO device_alerts (device_id,message,severity) VALUES ($1,$2,$3) RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, q, a.DeviceID, a.Message, a.Severity).Scan(&a.ID, &a.CreatedAt)
}

func (r *repository) ListAlerts(ctx context.Context, deviceID *string) ([]*Alert, error) {
	q := `SELECT da.id, da.device_id, d.name, da.message, da.severity, da.created_at, da.resolved_at
		FROM device_alerts da JOIN devices d ON d.id=da.device_id WHERE da.resolved_at IS NULL`
	args := []any{}
	if deviceID != nil {
		q += " AND da.device_id=$1"
		args = append(args, *deviceID)
	}
	q += " ORDER BY da.created_at DESC"
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var alerts []*Alert
	for rows.Next() {
		a := &Alert{}
		if err := rows.Scan(&a.ID, &a.DeviceID, &a.DeviceName, &a.Message, &a.Severity, &a.CreatedAt, &a.ResolvedAt); err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, rows.Err()
}

func (r *repository) ResolveAlert(ctx context.Context, alertID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE device_alerts SET resolved_at=NOW() WHERE id=$1`, alertID)
	return err
}

func (r *repository) AlertCount(ctx context.Context, deviceID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM device_alerts WHERE device_id=$1 AND resolved_at IS NULL`, deviceID).Scan(&count)
	return count, err
}
