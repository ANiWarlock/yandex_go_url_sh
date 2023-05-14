package main

import (
	"github.com/ANiWarlock/yandex_go_url_sh.git/internal/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))

	require.NoError(t, err)

	// вместо редиректа возвращаем предыдущий запрос
	client := ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func Test_Router(t *testing.T) {
	ts := httptest.NewServer(shortenerRouter())
	defer ts.Close()

	url := "http://ya.ru"
	shortURLHash := "shortened"
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
		{
			name:   "GET test #1",
			method: http.MethodGet,
			url:    "/" + shortURLHash,
			body:   "",
			want: want{
				code:        307,
				contentType: "",
				location:    url,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app.LinkStore = make(map[string]string)
			if test.method == http.MethodGet {
				app.LinkStore[shortURLHash] = url
			}

			resp, resBody := testRequest(t, ts, test.method, test.url, test.body)
			defer resp.Body.Close()
			
			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))

			if test.method == http.MethodPost {
				for key := range app.LinkStore {
					assert.Equal(t, key, resBody[len(resBody)-8:])
				}
			} else if test.method == http.MethodGet {
				assert.Equal(t, test.want.location, resp.Header.Get("Location"))
			}
		})
	}
}
