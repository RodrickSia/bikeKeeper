package vehicle

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, v *Vehicle) error
	GetByID(ctx context.Context, id string) (*Vehicle, error)
	ListByOwner(ctx context.Context, ownerID string) ([]*Vehicle, error)
	List(ctx context.Context) ([]*Vehicle, error)
	FindByPlate(ctx context.Context, plate string) (*Vehicle, error)
	Deactivate(ctx context.Context, id, ownerID string) error
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }

const cols = `id, license_plate, brand, model, color, owner_id, is_active, registered_at`

func scanVehicle(row *sql.Row, v *Vehicle) error {
	return row.Scan(&v.ID, &v.LicensePlate, &v.Brand, &v.Model, &v.Color, &v.OwnerID, &v.IsActive, &v.RegisteredAt)
}

func (r *repository) Create(ctx context.Context, v *Vehicle) error {
	const q = `INSERT INTO vehicles (license_plate,brand,model,color,owner_id) VALUES ($1,$2,$3,$4,$5) RETURNING ` + cols
	return scanVehicle(r.db.QueryRowContext(ctx, q, v.LicensePlate, v.Brand, v.Model, v.Color, v.OwnerID), v)
}

func (r *repository) GetByID(ctx context.Context, id string) (*Vehicle, error) {
	v := &Vehicle{}
	err := scanVehicle(r.db.QueryRowContext(ctx, `SELECT `+cols+` FROM vehicles WHERE id=$1`, id), v)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("vehicle not found")
	}
	return v, err
}

func (r *repository) FindByPlate(ctx context.Context, plate string) (*Vehicle, error) {
	v := &Vehicle{}
	err := scanVehicle(r.db.QueryRowContext(ctx, `SELECT `+cols+` FROM vehicles WHERE LOWER(REPLACE(REPLACE(license_plate, '-', ''), ' ', '')) = LOWER(REPLACE(REPLACE($1, '-', ''), ' ', '')) AND is_active=TRUE`, plate), v)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("vehicle not found")
	}
	return v, err
}

func (r *repository) ListByOwner(ctx context.Context, ownerID string) ([]*Vehicle, error) {
	return r.query(ctx, `SELECT `+cols+` FROM vehicles WHERE owner_id=$1 AND is_active=TRUE ORDER BY registered_at DESC`, ownerID)
}

func (r *repository) List(ctx context.Context) ([]*Vehicle, error) {
	return r.query(ctx, `SELECT `+cols+` FROM vehicles ORDER BY registered_at DESC`)
}

func (r *repository) query(ctx context.Context, q string, args ...any) ([]*Vehicle, error) {
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var vs []*Vehicle
	for rows.Next() {
		v := &Vehicle{}
		if err := rows.Scan(&v.ID, &v.LicensePlate, &v.Brand, &v.Model, &v.Color, &v.OwnerID, &v.IsActive, &v.RegisteredAt); err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}
	return vs, rows.Err()
}

func (r *repository) Deactivate(ctx context.Context, id, ownerID string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE vehicles SET is_active=FALSE WHERE id=$1 AND owner_id=$2`, id, ownerID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("vehicle %s not found", id)
	}
	return nil
}
