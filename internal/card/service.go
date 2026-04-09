package card

import (
	"context"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

type CreateParams struct {
	CardUID  string
	CardType string
	MemberID *string
}

func (s *Service) Create(ctx context.Context, params CreateParams) (*Card, error) {
	if params.CardType != "monthly" && params.CardType != "casual" {
		return nil, fmt.Errorf("invalid card type: %s", params.CardType)
	}

	card := &Card{
		CardUID:  params.CardUID,
		CardType: params.CardType,
		MemberID: params.MemberID,
		IsInside: false,
		Status:   "active",
	}

	if err := s.repo.Create(ctx, card); err != nil {
		return nil, fmt.Errorf("creating card: %w", err)
	}
	return card, nil
}

func (s *Service) GetByUID(ctx context.Context, cardUID string) (*Card, error) {
	return s.repo.GetByUID(ctx, cardUID)
}

func (s *Service) ListByMember(ctx context.Context, memberID string) ([]*Card, error) {
	return s.repo.ListByMember(ctx, memberID)
}

type UpdateParams struct {
	CardUID  string
	CardType *string
	MemberID *string
	Status   *string
}

func (s *Service) Update(ctx context.Context, params UpdateParams) (*Card, error) {
	card, err := s.repo.GetByUID(ctx, params.CardUID)
	if err != nil {
		return nil, err
	}

	if params.CardType != nil {
		if *params.CardType != "monthly" && *params.CardType != "casual" {
			return nil, fmt.Errorf("invalid card type: %s", *params.CardType)
		}
		card.CardType = *params.CardType
	}
	if params.MemberID != nil {
		card.MemberID = params.MemberID
	}
	if params.Status != nil {
		if *params.Status != "active" && *params.Status != "blocked" && *params.Status != "lost" {
			return nil, fmt.Errorf("invalid status: %s", *params.Status)
		}
		card.Status = *params.Status
	}

	if err := s.repo.Update(ctx, card); err != nil {
		return nil, fmt.Errorf("updating card: %w", err)
	}
	return card, nil
}

func (s *Service) ToggleInside(ctx context.Context, cardUID string) (*Card, error) {
	card, err := s.repo.GetByUID(ctx, cardUID)
	if err != nil {
		return nil, err
	}
	if card.Status != "active" {
		return nil, fmt.Errorf("card %s is %s", cardUID, card.Status)
	}
	card.IsInside = !card.IsInside
	if err := s.repo.Update(ctx, card); err != nil {
		return nil, fmt.Errorf("toggling card: %w", err)
	}
	return card, nil
}

func (s *Service) Delete(ctx context.Context, cardUID string) error {
	return s.repo.Delete(ctx, cardUID)
}
