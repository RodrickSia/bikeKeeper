package monthlypass

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, p CreateParams) (*MonthlyPass, error)
	ListByUser(ctx context.Context, userID string) ([]*MonthlyPass, error)
	ToggleAutoRenew(ctx context.Context, id string) (*MonthlyPass, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, p CreateParams) (*MonthlyPass, error) {
	m := &MonthlyPass{}
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO monthly_passes (user_id, vehicle_id, vehicle_plate, vehicle_brand, month, start_date, end_date, price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, vehicle_id, vehicle_plate, vehicle_brand, month,
		          start_date, end_date, price, status, is_auto_renew, purchased_at`,
		p.UserID, p.VehicleID, p.VehiclePlate, p.VehicleBrand,
		p.Month, p.StartDate, p.EndDate, p.Price,
	).Scan(&m.ID, &m.UserID, &m.VehicleID, &m.VehiclePlate, &m.VehicleBrand,
		&m.Month, &m.StartDate, &m.EndDate, &m.Price, &m.Status, &m.IsAutoRenew, &m.PurchasedAt)
	if err != nil {
		return nil, fmt.Errorf("creating monthly pass: %w", err)
	}
	return m, nil
}

func (r *repository) ListByUser(ctx context.Context, userID string) ([]*MonthlyPass, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, vehicle_id, vehicle_plate, vehicle_brand, month,
		       start_date, end_date, price, status, is_auto_renew, purchased_at
		FROM monthly_passes
		WHERE user_id = $1
		ORDER BY purchased_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passes []*MonthlyPass
	for rows.Next() {
		m := &MonthlyPass{}
		if err := rows.Scan(&m.ID, &m.UserID, &m.VehicleID, &m.VehiclePlate, &m.VehicleBrand,
			&m.Month, &m.StartDate, &m.EndDate, &m.Price, &m.Status, &m.IsAutoRenew, &m.PurchasedAt); err != nil {
			return nil, err
		}
		passes = append(passes, m)
	}
	return passes, rows.Err()
}

func (r *repository) ToggleAutoRenew(ctx context.Context, id string) (*MonthlyPass, error) {
	m := &MonthlyPass{}
	err := r.db.QueryRowContext(ctx, `
		UPDATE monthly_passes
		SET is_auto_renew = NOT is_auto_renew
		WHERE id = $1
		RETURNING id, user_id, vehicle_id, vehicle_plate, vehicle_brand, month,
		          start_date, end_date, price, status, is_auto_renew, purchased_at`, id,
	).Scan(&m.ID, &m.UserID, &m.VehicleID, &m.VehiclePlate, &m.VehicleBrand,
		&m.Month, &m.StartDate, &m.EndDate, &m.Price, &m.Status, &m.IsAutoRenew, &m.PurchasedAt)
	if err != nil {
		return nil, fmt.Errorf("toggling auto renew: %w", err)
	}
	return m, nil
}
