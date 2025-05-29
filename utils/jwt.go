package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("supersecret")

func GenerateToken(userID int, role string, depoID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"depo_id": depoID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
