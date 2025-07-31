package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTUser struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func BuildJWTClaims(user JWTUser, expireDuration time.Duration) jwt.MapClaims {
	return jwt.MapClaims{
		"userId":   user.UserID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      jwt.NewNumericDate(time.Now().Add(expireDuration)),
	}
}
