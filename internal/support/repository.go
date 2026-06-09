package support

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, t *Ticket) error
	GetByID(ctx context.Context, id string) (*Ticket, error)
	List(ctx context.Context, userID *string) ([]*Ticket, error)
	UpdateStatus(ctx context.Context, id, status string) error
	AddResponse(ctx context.Context, resp *Response) error
	GetResponses(ctx context.Context, ticketID string) ([]Response, error)
}

type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, t *Ticket) error {
	const q = `INSERT INTO support_tickets (user_id,category,subject,description) VALUES ($1,$2,$3,$4)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, q, t.UserID, t.Category, t.Subject, t.Description).
		Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *repository) GetByID(ctx context.Context, id string) (*Ticket, error) {
	t := &Ticket{}
	err := r.db.QueryRowContext(ctx, `SELECT id,user_id,category,subject,description,status,created_at,updated_at FROM support_tickets WHERE id=$1`, id).
		Scan(&t.ID, &t.UserID, &t.Category, &t.Subject, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ticket not found")
	}
	return t, err
}

func (r *repository) List(ctx context.Context, userID *string) ([]*Ticket, error) {
	q := `SELECT id,user_id,category,subject,description,status,created_at,updated_at FROM support_tickets`
	args := []any{}
	if userID != nil {
		q += " WHERE user_id=$1"
		args = append(args, *userID)
	}
	q += " ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tickets []*Ticket
	for rows.Next() {
		t := &Ticket{}
		if err := rows.Scan(&t.ID, &t.UserID, &t.Category, &t.Subject, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, rows.Err()
}

func (r *repository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE support_tickets SET status=$2, updated_at=NOW() WHERE id=$1`, id, status)
	return err
}

func (r *repository) AddResponse(ctx context.Context, resp *Response) error {
	const q = `INSERT INTO ticket_responses (ticket_id,sender_id,message,is_admin) VALUES ($1,$2,$3,$4)
		RETURNING id, created_at`
	err := r.db.QueryRowContext(ctx, q, resp.TicketID, resp.SenderID, resp.Message, resp.IsAdmin).
		Scan(&resp.ID, &resp.CreatedAt)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `UPDATE support_tickets SET updated_at=NOW() WHERE id=$1`, resp.TicketID)
	return err
}

func (r *repository) GetResponses(ctx context.Context, ticketID string) ([]Response, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,ticket_id,sender_id,message,is_admin,created_at FROM ticket_responses WHERE ticket_id=$1 ORDER BY created_at`, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resps []Response
	for rows.Next() {
		var resp Response
		if err := rows.Scan(&resp.ID, &resp.TicketID, &resp.SenderID, &resp.Message, &resp.IsAdmin, &resp.CreatedAt); err != nil {
			return nil, err
		}
		resps = append(resps, resp)
	}
	return resps, rows.Err()
}
