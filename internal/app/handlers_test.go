package app

import (
	"context"
	"github.com/ANiWarlock/yandex_go_url_sh.git/config"
	"github.com/ANiWarlock/yandex_go_url_sh.git/lib/auth"
	"github.com/ANiWarlock/yandex_go_url_sh.git/logger"
	"github.com/ANiWarlock/yandex_go_url_sh.git/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// тесты iter2, теперь проверяет только пост, остальное в main_test.go
func Test_GetShortURLHandler(t *testing.T) {

	sugar, err := logger.Initialize("info")
	require.NoError(t, err)
	cfg, err := config.InitConfig()
	require.NoError(t, err)
	ctx := context.Background()
	store, err := storage.InitStorage(ctx, *cfg)
	require.NoError(t, err)
	myApp := NewApp(cfg, store, sugar)

	url := "http://ya.ru"
	type want struct {
		code        int
		contentType string
		location    string
	}
	tests := []struct {
		name   string
		method string
		url    string
		body   string
		want   want
	}{
		{
			name:   "POST test #1",
			method: http.MethodPost,
			url:    "/",
			body:   url,
			want: want{
				code:        201,
				contentType: "text/plain; charset=utf-8",
				location:    "",
			},
		},
		{
			name:   "POST test #2",
			method: http.MethodPost,
			url:    "/",
			body:   "",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
				location:    "",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.url, strings.NewReader(test.body))
			rr := httptest.NewRecorder()
			setCookie(request)
			myApp.GetShortURLHandler(rr, request)
			res := rr.Result()
			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, test.want.contentType, rr.Header().Get("Content-Type"))

			if test.body != "" {
				resHashedURL := string(resBody[len(resBody)-8:])
				_, err := store.GetLongURL(request.Context(), resHashedURL)
				assert.NoError(t, err)
			}
		})
	}
}

func setCookie(r *http.Request) {
	userID := uuid.NewString()
	cookieStringValue := auth.BuildCookieStringValue(userID)

	newCookie := &http.Cookie{
		Name:     "auth",
		Value:    cookieStringValue,
		HttpOnly: true,
	}

	r.AddCookie(newCookie)
}
