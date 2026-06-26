package parkingsession

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type repository struct {
	db *sql.DB
}

type Repository interface {
	Create(ctx context.Context, session *ParkingSession) error
	GetByID(ctx context.Context, id int64) (*ParkingSession, error)
	GetOngoingSessionByCard(ctx context.Context, cardUID string) (*ParkingSession, error)
	GetOngoingByPlate(ctx context.Context, plate string) (*ParkingSession, error)
	ListByCard(ctx context.Context, cardUID string) ([]*ParkingSession, error)
	CheckOut(ctx context.Context, id int64, session *ParkingSession) error
	Delete(ctx context.Context, id int64) error
	GetVehiclePlateByCard(ctx context.Context, cardUID string) (string, error)
	IsCasualCard(ctx context.Context, cardUID string) (bool, error)
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, session *ParkingSession) error {
	const query = `
		INSERT INTO parking_sessions
			(card_uid, plate_in, img_plate_in_path, img_person_in_path, check_in_time)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id, check_in_time`

	return r.db.QueryRowContext(ctx, query,
		session.CardUID, session.PlateIn, session.ImgPlateInPath, session.ImgPersonInPath, time.Now(),
	).Scan(&session.ID, &session.CheckInTime)
}

func (r *repository) GetOngoingSessionByCard(ctx context.Context, cardUID string) (*ParkingSession, error) {
	const query = `
		SELECT id, card_uid, plate_in, img_plate_in_path, img_person_in_path,
		       check_in_time, plate_out, img_plate_out_path, img_person_out_path,
		       check_out_time, status
		FROM parking_sessions
		WHERE card_uid = $1 AND status = 'ongoing'
		LIMIT 1`

	session := &ParkingSession{}
	err := r.db.QueryRowContext(ctx, query, cardUID).Scan(
		&session.ID, &session.CardUID, &session.PlateIn, &session.ImgPlateInPath, &session.ImgPersonInPath,
		&session.CheckInTime, &session.PlateOut, &session.ImgPlateOutPath, &session.ImgPersonOutPath,
		&session.CheckOutTime, &session.Status,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return session, err
}

func (r *repository) GetOngoingByPlate(ctx context.Context, plate string) (*ParkingSession, error) {
	const query = `
		SELECT id, card_uid, plate_in, img_plate_in_path, img_person_in_path,
		       check_in_time, plate_out, img_plate_out_path, img_person_out_path,
		       check_out_time, status
		FROM parking_sessions
		WHERE REPLACE(plate_in, '-', '') ILIKE REPLACE($1, '-', '') AND status = 'ongoing'
		LIMIT 1`

	session := &ParkingSession{}
	err := r.db.QueryRowContext(ctx, query, plate).Scan(
		&session.ID, &session.CardUID, &session.PlateIn, &session.ImgPlateInPath, &session.ImgPersonInPath,
		&session.CheckInTime, &session.PlateOut, &session.ImgPlateOutPath, &session.ImgPersonOutPath,
		&session.CheckOutTime, &session.Status,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return session, err
}

func (r *repository) GetByID(ctx context.Context, id int64) (*ParkingSession, error) {
	const query = `
		SELECT id, card_uid, plate_in, img_plate_in_path, img_person_in_path,
		       check_in_time, plate_out, img_plate_out_path, img_person_out_path,
		       check_out_time, status
		FROM parking_sessions
		WHERE id = $1`

	session := &ParkingSession{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID, &session.CardUID, &session.PlateIn, &session.ImgPlateInPath, &session.ImgPersonInPath,
		&session.CheckInTime, &session.PlateOut, &session.ImgPlateOutPath, &session.ImgPersonOutPath,
		&session.CheckOutTime, &session.Status,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("parking session %d not found", id)
	}
	return session, err
}

func (r *repository) ListByCard(ctx context.Context, cardUID string) ([]*ParkingSession, error) {
	const query = `
		SELECT id, card_uid, plate_in, img_plate_in_path, img_person_in_path,
		       check_in_time, plate_out, img_plate_out_path, img_person_out_path,
		       check_out_time, status
		FROM parking_sessions
		WHERE card_uid = $1
		ORDER BY check_in_time DESC`

	rows, err := r.db.QueryContext(ctx, query, cardUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*ParkingSession
	for rows.Next() {
		session := &ParkingSession{}
		if err := rows.Scan(
			&session.ID, &session.CardUID, &session.PlateIn, &session.ImgPlateInPath, &session.ImgPersonInPath,
			&session.CheckInTime, &session.PlateOut, &session.ImgPlateOutPath, &session.ImgPersonOutPath,
			&session.CheckOutTime, &session.Status,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, rows.Err()
}

func (r *repository) CheckOut(ctx context.Context, id int64, session *ParkingSession) error {
	const query = `
		UPDATE parking_sessions
		SET
			plate_out           = $1,
			img_plate_out_path  = $2,
			img_person_out_path = $3,
			check_out_time      = $4,
			status              = 'completed'
		WHERE id = $5 AND status = 'ongoing'`

	result, err := r.db.ExecContext(ctx, query,
		session.PlateOut, session.ImgPlateOutPath, session.ImgPersonOutPath, time.Now(),
		id,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("session %d not found or already completed", id)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM parking_sessions WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("parking session %d not found", id)
	}
	return nil
}

func (r *repository) GetVehiclePlateByCard(ctx context.Context, cardUID string) (string, error) {
	const query = `
		SELECT v.license_plate 
		FROM vehicles v 
		JOIN users u ON v.owner_id = u.id 
		JOIN members m ON u.member_id = m.id 
		JOIN cards c ON c.member_id = m.id 
		WHERE c.card_uid = $1 AND v.is_active = true 
		LIMIT 1`
	var plate string
	err := r.db.QueryRowContext(ctx, query, cardUID).Scan(&plate)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return plate, err
}

func (r *repository) IsCasualCard(ctx context.Context, cardUID string) (bool, error) {
	const query = `SELECT card_type = 'casual' FROM cards WHERE card_uid = $1`
	var isCasual bool
	err := r.db.QueryRowContext(ctx, query, cardUID).Scan(&isCasual)
	if err == sql.ErrNoRows {
		return false, fmt.Errorf("card %s not found", cardUID)
	}
	return isCasual, err
}

