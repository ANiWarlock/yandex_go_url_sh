package auth

import (
	"errors"
	"fmt"
	"github.com/ANiWarlock/yandex_go_url_sh.git/config"
	"github.com/golang-jwt/jwt/v4"
	"log"
)

var ErrInvalidToken = errors.New("token is not valid")

var secret string

type ctxKey uint

const CtxKeyUserID ctxKey = iota

func SetSecretKey(cfg *config.AppConfig) {
	secret = cfg.SecretKey
	if secret == "" {
		secret = "mySecretKey"
	}
}

func GetUserID(tokenString string) (string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token string: %w", err)
	}
	if !token.Valid {
		log.Printf("Received invalid token: %s \n", tokenString)
		return "", ErrInvalidToken
	}

	userID := claims["userId"].(string)

	return userID, nil
}

func BuildCookieStringValue(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userID,
		})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed signing string: %w", err)
	}

	return tokenString, nil
}
