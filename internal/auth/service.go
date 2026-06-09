package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserFinder interface {
	GetByEmail(ctx context.Context, email string) (UserRecord, error)
}

type UserRecord struct {
	ID           string
	PasswordHash string
	Role         string
	MemberID     *string
	Status       string
}

type Service struct {
	repo      UserFinder
	jwtSecret []byte
}

func NewService(repo UserFinder) *Service {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "bikekeeper-dev-secret" // fallback for development
	}
	return &Service{
		repo:      repo,
		jwtSecret: []byte(secret),
	}
}

type UserResponse struct {
	ID       string  `json:"id"`
	Email    string  `json:"email"`
	Role     string  `json:"role"`
	MemberID *string `json:"memberId,omitempty"`
}

func (s *Service) Login(ctx context.Context, email, password string) (string, *UserResponse, error) {
	record, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, fmt.Errorf("user search failed: %w", err)
	}

	if record.Status != "active" {
		return "", nil, fmt.Errorf("account status is %s", record.Status)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(record.PasswordHash), []byte(password)); err != nil {
		return "", nil, fmt.Errorf("password mismatch: %w", err)
	}

	memberID := ""
	if record.MemberID != nil {
		memberID = *record.MemberID
	}

	claims := Claims{
		UserID:   record.ID,
		Role:     record.Role,
		MemberID: memberID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return tokenStr, &UserResponse{
		ID:       record.ID,
		Email:    email,
		Role:     record.Role,
		MemberID: record.MemberID,
	}, nil
}

func (s *Service) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
