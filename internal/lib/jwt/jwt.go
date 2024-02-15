package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"oliva-back/internal/models"
	"time"
)

func NewToken(user models.Users, duration time.Duration, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"username": user.Surname,
		"mail":     user.Email,
		"exp":      time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
