package user

import (
	"context"

	"github.com/RodrickSia/bikeKeeper/internal/auth"
)

type AuthAdapter struct {
	repo Repository
}

func NewAuthAdapter(repo Repository) *AuthAdapter {
	return &AuthAdapter{repo: repo}
}

func (a *AuthAdapter) GetByEmail(ctx context.Context, email string) (auth.UserRecord, error) {
	u, err := a.repo.GetByEmail(ctx, email)
	if err != nil {
		return auth.UserRecord{}, err
	}
	return auth.UserRecord{
		ID:           u.ID,
		PasswordHash: u.PasswordHash,
		Role:         u.Role,
		MemberID:     u.MemberID,
	}, nil
}
