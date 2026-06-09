package shift

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, s *Shift) error
	GetByID(ctx context.Context, id string) (*Shift, error)
	List(ctx context.Context, from, to, staffID *string) ([]*Shift, error)
	UpdateStatus(ctx context.Context, id, status string) error
	UpdateNotes(ctx context.Context, id string, notes *string) error
	Delete(ctx context.Context, id string) error
	AssignStaff(ctx context.Context, shiftID, userID string, assignedBy *string) error
	RemoveStaff(ctx context.Context, shiftID, userID string) error
	GetAssignments(ctx context.Context, shiftID string) ([]string, []string, error)
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, s *Shift) error {
	const q = `INSERT INTO shifts (name,type,start_time,end_time,date,notes) VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, q, s.Name, s.Type, s.StartTime, s.EndTime, s.Date, s.Notes).
		Scan(&s.ID, &s.CreatedAt)
}

func (r *repository) GetByID(ctx context.Context, id string) (*Shift, error) {
	s := &Shift{}
	err := r.db.QueryRowContext(ctx, `SELECT id,name,type,start_time,end_time,date,status,notes,created_at FROM shifts WHERE id=$1`, id).
		Scan(&s.ID, &s.Name, &s.Type, &s.StartTime, &s.EndTime, &s.Date, &s.Status, &s.Notes, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("shift not found")
	}
	return s, err
}

func (r *repository) List(ctx context.Context, from, to, staffID *string) ([]*Shift, error) {
	q := `SELECT DISTINCT s.id,s.name,s.type,s.start_time,s.end_time,s.date,s.status,s.notes,s.created_at FROM shifts s`
	args := []any{}
	i := 1
	if staffID != nil {
		q += ` JOIN shift_assignments sa ON sa.shift_id = s.id`
		args = append(args, *staffID)
		q += fmt.Sprintf(" WHERE sa.user_id = $%d", i)
		i++
	} else {
		q += " WHERE 1=1"
	}
	if from != nil {
		args = append(args, *from)
		q += fmt.Sprintf(" AND s.date >= $%d", i)
		i++
	}
	if to != nil {
		args = append(args, *to)
		q += fmt.Sprintf(" AND s.date <= $%d", i)
		i++
	}
	q += " ORDER BY s.date, s.start_time"
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	shifts := []*Shift{}
	for rows.Next() {
		s := &Shift{}
		if err := rows.Scan(&s.ID, &s.Name, &s.Type, &s.StartTime, &s.EndTime, &s.Date, &s.Status, &s.Notes, &s.CreatedAt); err != nil { // scan list
			return nil, err
		}
		shifts = append(shifts, s)
	}
	return shifts, rows.Err()
}

func (r *repository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE shifts SET status=$2 WHERE id=$1`, id, status)
	return err
}

func (r *repository) UpdateNotes(ctx context.Context, id string, notes *string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE shifts SET notes=$2 WHERE id=$1`, id, notes)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM shifts WHERE id=$1`, id)
	return err
}

func (r *repository) AssignStaff(ctx context.Context, shiftID, userID string, assignedBy *string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO shift_assignments (shift_id,user_id,assigned_by) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING`,
		shiftID, userID, assignedBy)
	return err
}

func (r *repository) RemoveStaff(ctx context.Context, shiftID, userID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM shift_assignments WHERE shift_id=$1 AND user_id=$2`, shiftID, userID)
	return err
}

func (r *repository) GetAssignments(ctx context.Context, shiftID string) ([]string, []string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT sa.user_id, u.email FROM shift_assignments sa
		JOIN users u ON u.id = sa.user_id WHERE sa.shift_id=$1`, shiftID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	ids := []string{}
	names := []string{}
	for rows.Next() {
		var id, email string
		if err := rows.Scan(&id, &email); err != nil {
			return nil, nil, err
		}
		ids = append(ids, id)
		names = append(names, email)
	}
	return ids, names, rows.Err()
}
