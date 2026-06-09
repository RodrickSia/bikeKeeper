package cardrequest

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, req *CardRequest) error
	GetByID(ctx context.Context, id string) (*CardRequest, error)
	ListByMember(ctx context.Context, memberID string) ([]*CardRequest, error)
	List(ctx context.Context, status *string) ([]*CardRequest, error)
	UpdateStatus(ctx context.Context, id, status string, cardUID, rejectedReason, reviewedBy *string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, req *CardRequest) error {
	const query = `
		INSERT INTO card_requests
			(member_id, vehicle_plate, vehicle_brand, vehicle_model, vehicle_color, id_card_number, note)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, submitted_at`
	return r.db.QueryRowContext(ctx, query,
		req.MemberID, req.VehiclePlate, req.VehicleBrand, req.VehicleModel,
		req.VehicleColor, req.IDCardNumber, req.Note,
	).Scan(&req.ID, &req.SubmittedAt)
}

func (r *repository) GetByID(ctx context.Context, id string) (*CardRequest, error) {
	const query = `
		SELECT id, member_id, vehicle_plate, vehicle_brand, vehicle_model, vehicle_color,
		       id_card_number, note, status, card_uid, rejected_reason,
		       submitted_at, reviewed_at, reviewed_by
		FROM card_requests WHERE id = $1`
	req := &CardRequest{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&req.ID, &req.MemberID, &req.VehiclePlate, &req.VehicleBrand, &req.VehicleModel,
		&req.VehicleColor, &req.IDCardNumber, &req.Note, &req.Status, &req.CardUID,
		&req.RejectedReason, &req.SubmittedAt, &req.ReviewedAt, &req.ReviewedBy,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("card_request %s not found", id)
	}
	return req, err
}

func (r *repository) ListByMember(ctx context.Context, memberID string) ([]*CardRequest, error) {
	const query = `
		SELECT id, member_id, vehicle_plate, vehicle_brand, vehicle_model, vehicle_color,
		       id_card_number, note, status, card_uid, rejected_reason,
		       submitted_at, reviewed_at, reviewed_by
		FROM card_requests WHERE member_id = $1 ORDER BY submitted_at DESC`
	return r.scan(ctx, query, memberID)
}

func (r *repository) List(ctx context.Context, status *string) ([]*CardRequest, error) {
	query := `
		SELECT id, member_id, vehicle_plate, vehicle_brand, vehicle_model, vehicle_color,
		       id_card_number, note, status, card_uid, rejected_reason,
		       submitted_at, reviewed_at, reviewed_by
		FROM card_requests`
	args := []any{}
	if status != nil {
		query += " WHERE status = $1"
		args = append(args, *status)
	}
	query += " ORDER BY submitted_at DESC"
	return r.scan(ctx, query, args...)
}

func (r *repository) scan(ctx context.Context, query string, args ...any) ([]*CardRequest, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []*CardRequest
	for rows.Next() {
		req := &CardRequest{}
		if err := rows.Scan(
			&req.ID, &req.MemberID, &req.VehiclePlate, &req.VehicleBrand, &req.VehicleModel,
			&req.VehicleColor, &req.IDCardNumber, &req.Note, &req.Status, &req.CardUID,
			&req.RejectedReason, &req.SubmittedAt, &req.ReviewedAt, &req.ReviewedBy,
		); err != nil {
			return nil, err
		}
		results = append(results, req)
	}
	return results, rows.Err()
}

func (r *repository) UpdateStatus(ctx context.Context, id, status string, cardUID, rejectedReason, reviewedBy *string) error {
	const query = `
		UPDATE card_requests
		SET status=$2, card_uid=$3, rejected_reason=$4, reviewed_by=$5, reviewed_at=NOW()
		WHERE id=$1`
	result, err := r.db.ExecContext(ctx, query, id, status, cardUID, rejectedReason, reviewedBy)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("card_request %s not found", id)
	}
	return nil
}
