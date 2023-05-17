package app

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func (a *App) GetShortURLHandler(rw http.ResponseWriter, r *http.Request) {
	responseData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Error in request body: %s", err), http.StatusBadRequest)
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

func shorten(longURL string) string {
	hashedString := sha1.Sum([]byte(longURL))
	return hex.EncodeToString(hashedString[:4])
}
