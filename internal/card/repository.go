package card

import (
	"context"
	"database/sql"
	"fmt"
)

type repository struct {
	db *sql.DB
}

type Repository interface {
	Create(ctx context.Context, card *Card) error
	GetByUID(ctx context.Context, cardUID string) (*Card, error)
	ListByMember(ctx context.Context, memberID string) ([]*Card, error)
	GetAvailableCasual(ctx context.Context) (*Card, error)
	Update(ctx context.Context, card *Card) error
	Delete(ctx context.Context, cardUID string) error
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, card *Card) error {
	const query = `
		INSERT INTO cards (card_uid, card_type, member_id, is_inside, status, balance)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		card.CardUID, card.CardType, card.MemberID, card.IsInside, card.Status, card.Balance,
	)
	return err
}

func (r *repository) GetByUID(ctx context.Context, cardUID string) (*Card, error) {
	const query = `
		SELECT card_uid, card_type, member_id, is_inside, status, balance
		FROM cards
		WHERE card_uid = $1`

	card := &Card{}
	err := r.db.QueryRowContext(ctx, query, cardUID).Scan(
		&card.CardUID, &card.CardType, &card.MemberID, &card.IsInside, &card.Status, &card.Balance,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("card %s not found", cardUID)
	}
	return card, err
}

func (r *repository) ListByMember(ctx context.Context, memberID string) ([]*Card, error) {
	const query = `
		SELECT card_uid, card_type, member_id, is_inside, status, balance
		FROM cards
		WHERE member_id = $1`

	rows, err := r.db.QueryContext(ctx, query, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*Card
	for rows.Next() {
		card := &Card{}
		if err := rows.Scan(
			&card.CardUID, &card.CardType, &card.MemberID, &card.IsInside, &card.Status, &card.Balance,
		); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, rows.Err()
}

func (r *repository) GetAvailableCasual(ctx context.Context) (*Card, error) {
	const query = `
		SELECT c.card_uid, c.card_type, c.member_id, c.is_inside, c.status, c.balance
		FROM cards c
		LEFT JOIN parking_sessions ps ON ps.card_uid = c.card_uid AND ps.status = 'ongoing'
		WHERE c.card_type = 'casual' AND c.status = 'active' AND ps.id IS NULL
		LIMIT 1`

	card := &Card{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&card.CardUID, &card.CardType, &card.MemberID, &card.IsInside, &card.Status, &card.Balance,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return card, err
}

func (r *repository) Update(ctx context.Context, card *Card) error {
	const query = `
		UPDATE cards
		SET card_type = $1, member_id = $2, is_inside = $3, status = $4, balance = $5
		WHERE card_uid = $6`

	result, err := r.db.ExecContext(ctx, query,
		card.CardType, card.MemberID, card.IsInside, card.Status, card.Balance, card.CardUID,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("card %s not found", card.CardUID)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, cardUID string) error {
	const query = `DELETE FROM cards WHERE card_uid = $1`
	result, err := r.db.ExecContext(ctx, query, cardUID)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("card %s not found", cardUID)
	}
	return nil
}
