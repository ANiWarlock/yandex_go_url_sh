package app

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/ANiWarlock/yandex_go_url_sh.git/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type APIRequest struct {
	URL string `json:"url"`
}

type APIResponse struct {
	Result string `json:"result"`
}

func (a *App) GetShortURLHandler(rw http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var apiReq APIRequest
	var apiResp APIResponse
	var longURL string
	var resp []byte

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		a.sugar.Errorf("Cannot process body: %v", err)
		return
	}

	switch r.URL.Path {
	case "/":
		if buf.String() == "" {
			http.Error(rw, "Empty body!", http.StatusBadRequest)
			return
		}
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		longURL = buf.String()
	case "/api/shorten":
		if err = json.Unmarshal(buf.Bytes(), &apiReq); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			a.sugar.Errorf("Cannot process body: %v", err)
			return
		}

		if apiReq.URL == "" {
			http.Error(rw, "Empty URL!", http.StatusBadRequest)
			return
		}
		rw.Header().Set("Content-Type", "application/json")

		longURL = apiReq.URL
	}

	hashedURL := shorten(longURL)
	shortURL := a.cfg.BaseURL + "/" + hashedURL

	if err = a.storage.SaveLongURL(hashedURL, longURL); err != nil {
		if errors.Unwrap(err) == storage.ErrUniqueViolation {

			switch r.URL.Path {
			case "/":
				resp = []byte(shortURL)
			case "/api/shorten":
				apiResp.Result = shortURL

				resp, err = json.Marshal(apiResp)
				if err != nil {
					rw.WriteHeader(http.StatusBadRequest)
					a.sugar.Errorf("Cannot process body: %v", err)
					return
				}
			}

			rw.WriteHeader(http.StatusConflict)
			_, err = rw.Write(resp)
			if err != nil {
				return
			}
			return
		}
		http.Error(rw, "Server Error", http.StatusBadRequest)
	}

	switch r.URL.Path {
	case "/":
		resp = []byte(shortURL)
	case "/api/shorten":
		apiResp.Result = shortURL

		resp, err = json.Marshal(apiResp)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			a.sugar.Errorf("Cannot process body: %v", err)
			return
		}
	}

	rw.WriteHeader(http.StatusCreated)
	_, err = rw.Write(resp)
	if err != nil {
		return
	}

}

func (a *App) LongURLRedirectHandler(rw http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")

	if shortURL == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	item, err := a.storage.GetLongURL(shortURL)
	if err != nil {
		http.Error(rw, "Not Found", http.StatusNotFound)
		return
	}

	rw.Header().Set("Location", item.LongURL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

type APIBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type APIBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (a *App) APIBatchHandler(rw http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var apiReq []APIBatchRequest
	var apiResp []APIBatchResponse
	var items []storage.Item

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		a.sugar.Errorf("Cannot process body: %v", err)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &apiReq); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		a.sugar.Errorf("Cannot process body: %v", err)
		return
	}

	// собираем ответ
	for _, v := range apiReq {
		if v.OriginalURL == "" {
			http.Error(rw, "Empty URL!", http.StatusBadRequest)
			return
		}

		if v.CorrelationID == "" {
			http.Error(rw, "Empty correlation ID!", http.StatusBadRequest)
			return
		}

		longURL := v.OriginalURL
		hashedURL := shorten(longURL)
		shortURL := a.cfg.BaseURL + "/" + hashedURL

		items = append(items, storage.Item{LongURL: longURL, ShortURL: hashedURL})
		apiResp = append(apiResp, APIBatchResponse{CorrelationID: v.CorrelationID, ShortURL: shortURL})
	}

	// сохраняем
	err = a.storage.BatchInsert(items)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		a.sugar.Errorf("Cannot batch save: %v", err)
		return
	}

	resp, err := json.Marshal(apiResp)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		a.sugar.Errorf("Cannot process body: %v", err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	_, err = rw.Write(resp)
	if err != nil {
		return
	}
}

func (a *App) PingHandler(rw http.ResponseWriter, r *http.Request) {
	if err := a.storage.Ping(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		rw.WriteHeader(http.StatusOK)
	}
}

func shorten(longURL string) string {
	hashedString := sha1.Sum([]byte(longURL))
	return hex.EncodeToString(hashedString[:4])
}
