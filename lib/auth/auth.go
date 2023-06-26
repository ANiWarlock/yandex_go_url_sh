package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

var ErrInvalidToken = errors.New("token is not valid")

const secret = "secretKey"

func GetUserID(tokenString string) (string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if !token.Valid {
		return "", ErrInvalidToken
	}
	if err != nil {
		return "", fmt.Errorf("jwt parse failed: %w", err)
	}

	userID := claims["userId"].(string)

	return userID, nil
}

func BuildCookieStringValue(userID string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userID,
		})
	tokenString, _ := token.SignedString([]byte(secret))

	return tokenString
}
