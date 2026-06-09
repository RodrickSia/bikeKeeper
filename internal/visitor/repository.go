package visitor

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, p *VisitorPass) error
	GetByID(ctx context.Context, id string) (*VisitorPass, error)
	ListByUser(ctx context.Context, userID string) ([]*VisitorPass, error)
	UpdateStatus(ctx context.Context, id, status string) error
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }

const cols = `id, user_id, visitor_name, visitor_phone, vehicle_plate, valid_date, status, qr_code_data, created_at, used_at`

func scanRow(row *sql.Row, p *VisitorPass) error {
	var validDate sql.NullString
	err := row.Scan(&p.ID, &p.UserID, &p.VisitorName, &p.VisitorPhone, &p.VehiclePlate, &validDate, &p.Status, &p.QRCodeData, &p.CreatedAt, &p.UsedAt)
	if err != nil {
		return err
	}
	p.ValidDate = validDate.String
	return nil
}

func (r *repository) Create(ctx context.Context, p *VisitorPass) error {
	const q = `INSERT INTO visitor_passes (user_id,visitor_name,visitor_phone,vehicle_plate,valid_date,qr_code_data)
		VALUES ($1,$2,$3,$4,$5,$6) RETURNING ` + cols
	return scanRow(r.db.QueryRowContext(ctx, q, p.UserID, p.VisitorName, p.VisitorPhone, p.VehiclePlate, p.ValidDate, p.QRCodeData), p)
}

func (r *repository) GetByID(ctx context.Context, id string) (*VisitorPass, error) {
	p := &VisitorPass{}
	err := scanRow(r.db.QueryRowContext(ctx, `SELECT `+cols+` FROM visitor_passes WHERE id=$1`, id), p)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("visitor pass not found")
	}
	return p, err
}

func (r *repository) ListByUser(ctx context.Context, userID string) ([]*VisitorPass, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+cols+` FROM visitor_passes WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var passes []*VisitorPass
	for rows.Next() {
		p := &VisitorPass{}
		var validDate sql.NullString
		if err := rows.Scan(&p.ID, &p.UserID, &p.VisitorName, &p.VisitorPhone, &p.VehiclePlate, &validDate, &p.Status, &p.QRCodeData, &p.CreatedAt, &p.UsedAt); err != nil {
			return nil, err
		}
		p.ValidDate = validDate.String
		passes = append(passes, p)
	}
	return passes, rows.Err()
}

func (r *repository) UpdateStatus(ctx context.Context, id, status string) error {
	q := `UPDATE visitor_passes SET status=$2`
	if status == "used" {
		q += `, used_at=NOW()`
	}
	q += ` WHERE id=$1`
	_, err := r.db.ExecContext(ctx, q, id, status)
	return err
}
