package ds

import (
	"development-of-internet-applications/internal/app/role"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID uint64    `json:"user_id"`
	Role   role.Role `json:"role"`
}