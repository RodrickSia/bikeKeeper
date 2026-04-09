package member

import (
	"context"
	"database/sql"
	"fmt"
)

type repository struct {
	db *sql.DB
}

type Repository interface {
	Create(ctx context.Context, member *Member) error
	GetByID(ctx context.Context, id string) (*Member, error)
	GetByStudentID(ctx context.Context, studentID string) (*Member, error)
	List(ctx context.Context) ([]*Member, error)
	Update(ctx context.Context, member *Member) error
	Delete(ctx context.Context, id string) error
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, member *Member) error {
	const query = `
		INSERT INTO members (student_id, full_name, phone)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	return r.db.QueryRowContext(ctx, query,
		member.StudentID, member.FullName, member.Phone,
	).Scan(&member.ID, &member.CreatedAt)
}

func (r *repository) GetByID(ctx context.Context, id string) (*Member, error) {
	const query = `
		SELECT id, student_id, full_name, phone, created_at
		FROM members
		WHERE id = $1`

	member := &Member{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&member.ID, &member.StudentID, &member.FullName, &member.Phone, &member.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("member %s not found", id)
	}
	return member, err
}

func (r *repository) GetByStudentID(ctx context.Context, studentID string) (*Member, error) {
	const query = `
		SELECT id, student_id, full_name, phone, created_at
		FROM members
		WHERE student_id = $1`

	member := &Member{}
	err := r.db.QueryRowContext(ctx, query, studentID).Scan(
		&member.ID, &member.StudentID, &member.FullName, &member.Phone, &member.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("member with student_id %s not found", studentID)
	}
	return member, err
}

func (r *repository) List(ctx context.Context) ([]*Member, error) {
	const query = `
		SELECT id, student_id, full_name, phone, created_at
		FROM members
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*Member
	for rows.Next() {
		member := &Member{}
		if err := rows.Scan(
			&member.ID, &member.StudentID, &member.FullName, &member.Phone, &member.CreatedAt,
		); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, rows.Err()
}

func (r *repository) Update(ctx context.Context, member *Member) error {
	const query = `
		UPDATE members
		SET full_name = $1, phone = $2
		WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query,
		member.FullName, member.Phone, member.ID,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("member %s not found", member.ID)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM members WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("member %s not found", id)
	}
	return nil
}
