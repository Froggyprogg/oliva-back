package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	TokenExpireDuration = time.Hour * 24 * 7
)

func GenerateToken(login string, tokenID uint) (string, error) {
	privateKey, err := PrivateKeyToRsa()
	if err != nil {
		return "", err
	}

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   login,
		Issuer:    "oliva-sso",
		ID:        string(tokenID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
