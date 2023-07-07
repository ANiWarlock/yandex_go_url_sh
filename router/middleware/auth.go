package middleware

import (
	"context"
	"fmt"
	"github.com/ANiWarlock/yandex_go_url_sh.git/lib/auth"
	"github.com/google/uuid"
	"net/http"
)

func SetAuthCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		tokenString, err := r.Cookie("auth")
		if err != nil {
			userID, err := setCookie(rw, r)
			if err != nil {
				http.Error(rw, "500", http.StatusInternalServerError)
				return
			}
			r = setCtxUserID(r, userID)
			next.ServeHTTP(rw, r)
			return
		}

		userID, err := auth.GetUserID(tokenString.Value)
		if err != nil {
			userID, err := setCookie(rw, r)
			if err != nil {
				http.Error(rw, "500", http.StatusInternalServerError)
				return
			}
			r = setCtxUserID(r, userID)
		}
		if userID == "" {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r = setCtxUserID(r, userID)
		next.ServeHTTP(rw, r)
	})
}

func setCookie(rw http.ResponseWriter, r *http.Request) (string, error) {
	userID := uuid.NewString()
	cookieStringValue, err := auth.BuildCookieStringValue(userID)
	if err != nil {
		return "", fmt.Errorf("failed to build cookie string: %w", err)
	}

	newCookie := &http.Cookie{
		Name:     "auth",
		Value:    cookieStringValue,
		HttpOnly: true,
	}

	http.SetCookie(rw, newCookie)
	return userID, nil
}

func setCtxUserID(r *http.Request, userID string) *http.Request {
	ctx := context.WithValue(r.Context(), auth.CtxKeyUserID, userID)
	return r.WithContext(ctx)
}
