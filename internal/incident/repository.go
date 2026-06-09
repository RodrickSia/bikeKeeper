package incident

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, inc *Incident) error
	GetByID(ctx context.Context, id string) (*Incident, error)
	List(ctx context.Context, staffID *string) ([]*Incident, error)
	Resolve(ctx context.Context, id, note string) error
	Escalate(ctx context.Context, id string) error
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }

const cols = `id, reported_by, vehicle_plate, type, description, location, status, resolved_at, resolved_note, created_at`

func scanRow(row *sql.Row, inc *Incident) error {
	return row.Scan(&inc.ID, &inc.ReportedBy, &inc.VehiclePlate, &inc.Type, &inc.Description,
		&inc.Location, &inc.Status, &inc.ResolvedAt, &inc.ResolvedNote, &inc.CreatedAt)
}

func (r *repository) Create(ctx context.Context, inc *Incident) error {
	const q = `INSERT INTO incidents (reported_by,vehicle_plate,type,description,location)
		VALUES ($1,$2,$3,$4,$5) RETURNING ` + cols
	return scanRow(r.db.QueryRowContext(ctx, q, inc.ReportedBy, inc.VehiclePlate, inc.Type, inc.Description, inc.Location), inc)
}

func (r *repository) GetByID(ctx context.Context, id string) (*Incident, error) {
	inc := &Incident{}
	err := scanRow(r.db.QueryRowContext(ctx, `SELECT `+cols+` FROM incidents WHERE id=$1`, id), inc)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("incident not found")
	}
	return inc, err
}

func (r *repository) List(ctx context.Context, staffID *string) ([]*Incident, error) {
	q := `SELECT ` + cols + ` FROM incidents`
	args := []any{}
	if staffID != nil {
		q += " WHERE reported_by=$1"
		args = append(args, *staffID)
	}
	q += " ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var incs []*Incident
	for rows.Next() {
		inc := &Incident{}
		if err := rows.Scan(&inc.ID, &inc.ReportedBy, &inc.VehiclePlate, &inc.Type, &inc.Description,
			&inc.Location, &inc.Status, &inc.ResolvedAt, &inc.ResolvedNote, &inc.CreatedAt); err != nil {
			return nil, err
		}
		incs = append(incs, inc)
	}
	return incs, rows.Err()
}

func (r *repository) Resolve(ctx context.Context, id, note string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE incidents SET status='resolved', resolved_at=NOW(), resolved_note=$2 WHERE id=$1`, id, note)
	return err
}

func (r *repository) Escalate(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE incidents SET status='escalated' WHERE id=$1`, id)
	return err
}
