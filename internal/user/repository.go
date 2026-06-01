package user

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	const query = `
		INSERT INTO users (email, password_hash, role, member_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	return r.db.QueryRowContext(ctx, query,
		user.Email, user.PasswordHash, user.Role, user.MemberID,
	).Scan(&user.ID, &user.CreatedAt)
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	const query = `
		SELECT id, email, password_hash, role, member_id, created_at
		FROM users
		WHERE email = $1`

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.MemberID, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return user, err
}

func (r *repository) GetByID(ctx context.Context, id string) (*User, error) {
	const query = `
		SELECT id, email, password_hash, role, member_id, created_at
		FROM users
		WHERE id = $1`

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.MemberID, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user %s not found", id)
	}
	return user, err
}

func (r *repository) List(ctx context.Context) ([]*User, error) {
	const query = `
		SELECT id, email, password_hash, role, member_id, created_at
		FROM users
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		u := &User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.MemberID, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *repository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM users WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user %s not found", id)
	}
	return nil
}
