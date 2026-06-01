package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   string `json:"userId"`
	Role     string `json:"role"`
	MemberID string `json:"memberId,omitempty"`
	jwt.RegisteredClaims
}
