package app

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func (a *App) GetShortURLHandler(rw http.ResponseWriter, r *http.Request) {
	responseData, err := io.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		a.sugar.Errorf("Cannot process body: %v", err)
		return
	}
	if string(responseData) == "" {
		http.Error(rw, "Empty body!", http.StatusBadRequest)
		return
	}

	longURL := string(responseData)
	hashedURL := shorten(longURL)
	a.storage.SaveLongURL(hashedURL, longURL)
	shortURL := a.cfg.BaseURL + "/" + hashedURL

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err = rw.Write([]byte(shortURL))
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

	longURL, ok := a.storage.GetLongURL(shortURL)
	if !ok {
		http.Error(rw, "Not Found", http.StatusNotFound)
		return
	}

	rw.Header().Set("Location", longURL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

type APIRequest struct {
	URL string `json:"url"`
}

type APIResponse struct {
	Result string `json:"result"`
}

func (a *App) APIGetShortURLHandler(rw http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var apiReq APIRequest
	var apiResp APIResponse

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

	if apiReq.URL == "" {
		http.Error(rw, "Empty URL!", http.StatusBadRequest)
		return
	}

	longURL := apiReq.URL
	hashedURL := shorten(longURL)
	a.storage.SaveLongURL(hashedURL, longURL)
	shortURL := a.cfg.BaseURL + "/" + hashedURL
	apiResp.Result = shortURL

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

func shorten(longURL string) string {
	hashedString := sha1.Sum([]byte(longURL))
	return hex.EncodeToString(hashedString[:4])
}
