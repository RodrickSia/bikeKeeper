package parkinglot

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, lot *ParkingLot) error
	GetByID(ctx context.Context, id string) (*ParkingLot, error)
	List(ctx context.Context) ([]*ParkingLot, error)
	Update(ctx context.Context, lot *ParkingLot) error
	Delete(ctx context.Context, id string) error
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }

func scanLot(row *sql.Row, l *ParkingLot) error {
	return row.Scan(&l.ID, &l.Name, &l.Address, &l.Type, &l.Status,
		&l.TotalCapacity, &l.CurrentOccupancy, &l.OpenTime, &l.CloseTime,
		&l.ContactPhone, &l.ManagerName, &l.Description, &l.CreatedAt, &l.UpdatedAt)
}

const selectCols = `id, name, address, type, status, total_capacity, current_occupancy,
	open_time, close_time, contact_phone, manager_name, description, created_at, updated_at`

func (r *repository) Create(ctx context.Context, l *ParkingLot) error {
	const q = `INSERT INTO parking_lots (name,address,type,status,total_capacity,open_time,close_time,contact_phone,manager_name,description)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING ` + selectCols
	return scanLot(r.db.QueryRowContext(ctx, q,
		l.Name, l.Address, l.Type, l.Status, l.TotalCapacity,
		l.OpenTime, l.CloseTime, l.ContactPhone, l.ManagerName, l.Description,
	), l)
}

func (r *repository) GetByID(ctx context.Context, id string) (*ParkingLot, error) {
	l := &ParkingLot{}
	err := scanLot(r.db.QueryRowContext(ctx, `SELECT `+selectCols+` FROM parking_lots WHERE id=$1`, id), l)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("parking lot not found")
	}
	return l, err
}

func (r *repository) List(ctx context.Context) ([]*ParkingLot, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+selectCols+` FROM parking_lots ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var lots []*ParkingLot
	for rows.Next() {
		l := &ParkingLot{}
		if err := rows.Scan(&l.ID, &l.Name, &l.Address, &l.Type, &l.Status,
			&l.TotalCapacity, &l.CurrentOccupancy, &l.OpenTime, &l.CloseTime,
			&l.ContactPhone, &l.ManagerName, &l.Description, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		lots = append(lots, l)
	}
	return lots, rows.Err()
}

func (r *repository) Update(ctx context.Context, l *ParkingLot) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE parking_lots SET name=$2,address=$3,type=$4,status=$5,total_capacity=$6,
		current_occupancy=$7,open_time=$8,close_time=$9,contact_phone=$10,manager_name=$11,
		description=$12,updated_at=NOW() WHERE id=$1`,
		l.ID, l.Name, l.Address, l.Type, l.Status, l.TotalCapacity, l.CurrentOccupancy,
		l.OpenTime, l.CloseTime, l.ContactPhone, l.ManagerName, l.Description)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM parking_lots WHERE id=$1`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("parking lot %s not found", id)
	}
	return nil
}
