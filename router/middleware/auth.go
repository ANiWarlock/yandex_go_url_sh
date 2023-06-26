package middleware

import (
	"github.com/ANiWarlock/yandex_go_url_sh.git/lib/auth"
	"github.com/google/uuid"
	"net/http"
)

func SetAuthCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		tokenString, err := r.Cookie("auth")
		if err != nil {
			setCookie(rw, r)
			next.ServeHTTP(rw, r)
			return
		}

		userID, err := auth.GetUserID(tokenString.Value)
		if err != nil {
			setCookie(rw, r)
		}
		if userID == "" {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(rw, r)
	})
}

func setCookie(rw http.ResponseWriter, r *http.Request) {
	userID := uuid.NewString()
	cookieStringValue := auth.BuildCookieStringValue(userID)

	newCookie := &http.Cookie{
		Name:     "auth",
		Value:    cookieStringValue,
		HttpOnly: true,
	}

	r.AddCookie(newCookie)
	http.SetCookie(rw, newCookie)
}
