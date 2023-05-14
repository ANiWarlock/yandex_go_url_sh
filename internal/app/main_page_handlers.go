package app

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

const baseURL = "http://localhost:8080"

var LinkStore map[string]string

func MainPageHandler(rw http.ResponseWriter, r *http.Request) {
	LinkStore = make(map[string]string)

	defer r.Body.Close()
	responseData, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, fmt.Sprintf("Error in request body: %s", err), http.StatusBadRequest)
	}
	if string(responseData) == "" {
		http.Error(rw, "Empty body!", http.StatusBadRequest)
		return
	}

	longURL := string(responseData)
	hashedURL := shorten(longURL)
	LinkStore[hashedURL] = longURL
	shortURL := baseURL + "/" + hashedURL

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err = rw.Write([]byte(shortURL))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func LongUrlRedirectHandler(rw http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")

	fmt.Println(shortURL)
	fmt.Println(LinkStore[shortURL])
	fmt.Println("I WAS HERE")

	if shortURL == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.Header().Set("Location", LinkStore[shortURL])
	rw.WriteHeader(http.StatusTemporaryRedirect)
	return
}

func shorten(longURL string) string {
	hashedString := sha1.Sum([]byte(longURL))
	return hex.EncodeToString(hashedString[:4])
}
