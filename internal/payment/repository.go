package payment

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Deposit(ctx context.Context, cardUID string, amount float64) (*Transaction, error)
	Withdraw(ctx context.Context, cardUID string, amount float64, txType string) (*Transaction, error)
	ChargeParking(ctx context.Context, cardUID string, fee float64, sessionID int64) (*Transaction, error)
	ListByCard(ctx context.Context, cardUID string) ([]*Transaction, error)
	GetBalance(ctx context.Context, cardUID string) (float64, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Deposit(ctx context.Context, cardUID string, amount float64) (*Transaction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx,
		`UPDATE cards SET balance = balance + $1 WHERE card_uid = $2`,
		amount, cardUID,
	)
	if err != nil {
		return nil, fmt.Errorf("updating balance: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("checking update result: %w", err)
	}
	if affected == 0 {
		return nil, fmt.Errorf("card %s not found", cardUID)
	}

	t := &Transaction{}
	err = tx.QueryRowContext(ctx,
		`INSERT INTO transactions (card_uid, amount, type)
		 VALUES ($1, $2, $3)
		 RETURNING id, card_uid, amount, type, payment_method, session_id, created_at`,
		cardUID, amount, TypeDeposit,
	).Scan(&t.ID, &t.CardUID, &t.Amount, &t.Type, &t.PaymentMethod, &t.SessionID, &t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("inserting transaction: %w", err)
	}

	return t, tx.Commit()
}

func (r *repository) Withdraw(ctx context.Context, cardUID string, amount float64, txType string) (*Transaction, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("withdraw amount must be positive")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var balance float64
	err = tx.QueryRowContext(ctx,
		`SELECT balance FROM cards WHERE card_uid = $1 FOR UPDATE`,
		cardUID,
	).Scan(&balance)
	if err != nil {
		return nil, fmt.Errorf("fetching balance: %w", err)
	}
	if balance < amount {
		return nil, fmt.Errorf("insufficient balance: have %.0f, need %.0f", balance, amount)
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE cards SET balance = balance - $1 WHERE card_uid = $2`,
		amount, cardUID,
	)
	if err != nil {
		return nil, fmt.Errorf("deducting balance: %w", err)
	}

	t := &Transaction{}
	err = tx.QueryRowContext(ctx,
		`INSERT INTO transactions (card_uid, amount, type)
		 VALUES ($1, $2, $3)
		 RETURNING id, card_uid, amount, type, payment_method, session_id, created_at`,
		cardUID, amount, txType,
	).Scan(&t.ID, &t.CardUID, &t.Amount, &t.Type, &t.PaymentMethod, &t.SessionID, &t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("inserting transaction: %w", err)
	}

	return t, tx.Commit()
}

func (r *repository) ChargeParking(ctx context.Context, cardUID string, fee float64, sessionID int64) (*Transaction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var balance float64
	err = tx.QueryRowContext(ctx,
		`SELECT balance FROM cards WHERE card_uid = $1 FOR UPDATE`,
		cardUID,
	).Scan(&balance)
	if err != nil {
		return nil, fmt.Errorf("fetching balance: %w", err)
	}

	method := MethodCash
	if balance >= fee {
		method = MethodCardBalance
		_, err = tx.ExecContext(ctx,
			`UPDATE cards SET balance = balance - $1 WHERE card_uid = $2`,
			fee, cardUID,
		)
		if err != nil {
			return nil, fmt.Errorf("deducting balance: %w", err)
		}
	}

	t := &Transaction{}
	err = tx.QueryRowContext(ctx,
		`INSERT INTO transactions (card_uid, amount, type, payment_method, session_id)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, card_uid, amount, type, payment_method, session_id, created_at`,
		cardUID, fee, TypeParkingFee, method, sessionID,
	).Scan(&t.ID, &t.CardUID, &t.Amount, &t.Type, &t.PaymentMethod, &t.SessionID, &t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("inserting transaction: %w", err)
	}

	return t, tx.Commit()
}

func (r *repository) ListByCard(ctx context.Context, cardUID string) ([]*Transaction, error) {
	const query = `
		SELECT id, card_uid, amount, type, payment_method, session_id, created_at
		FROM transactions
		WHERE card_uid = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, cardUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []*Transaction
	for rows.Next() {
		t := &Transaction{}
		if err := rows.Scan(&t.ID, &t.CardUID, &t.Amount, &t.Type, &t.PaymentMethod, &t.SessionID, &t.CreatedAt); err != nil {
			return nil, err
		}
		txns = append(txns, t)
	}
	return txns, rows.Err()
}

func (r *repository) GetBalance(ctx context.Context, cardUID string) (float64, error) {
	var balance float64
	err := r.db.QueryRowContext(ctx,
		`SELECT balance FROM cards WHERE card_uid = $1`, cardUID,
	).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("card %s not found", cardUID)
	}
	return balance, err
}
