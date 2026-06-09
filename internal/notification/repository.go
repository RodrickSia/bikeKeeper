package notification

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, n *Notification) error
	ListByUser(ctx context.Context, userID string) ([]*Notification, error)
	MarkRead(ctx context.Context, id string) error
	MarkAllRead(ctx context.Context, userID string) error
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }

const cols = `id, user_id, title, message, type, is_read, created_at`

func (r *repository) Create(ctx context.Context, n *Notification) error {
	const q = `INSERT INTO notifications (user_id,title,message,type) VALUES ($1,$2,$3,$4) RETURNING ` + cols
	return r.db.QueryRowContext(ctx, q, n.UserID, n.Title, n.Message, n.Type).
		Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.Type, &n.IsRead, &n.CreatedAt)
}

func (r *repository) ListByUser(ctx context.Context, userID string) ([]*Notification, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+cols+` FROM notifications WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ns []*Notification
	for rows.Next() {
		n := &Notification{}
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.Type, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		ns = append(ns, n)
	}
	return ns, rows.Err()
}

func (r *repository) MarkRead(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE notifications SET is_read=TRUE WHERE id=$1`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("notification %s not found", id)
	}
	return nil
}

func (r *repository) MarkAllRead(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE notifications SET is_read=TRUE WHERE user_id=$1`, userID)
	return err
}
