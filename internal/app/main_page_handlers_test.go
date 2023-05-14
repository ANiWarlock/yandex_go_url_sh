package app

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// тесты iter2, теперь проверяет только пост, остальное в main_test.go
func Test_mainPageHandler(t *testing.T) {
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
		t.Run(test.name, func(t *testing.T) {
			LinkStore = make(map[string]string)

			request := httptest.NewRequest(test.method, test.url, strings.NewReader(test.body))
			rr := httptest.NewRecorder()
			MainPageHandler(rr, request)
			res := rr.Result()
			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, test.want.contentType, rr.Header().Get("Content-Type"))

			for key := range LinkStore {
				assert.Equal(t, baseURL+"/"+key, string(resBody))
			}
		})
	}
}
