package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

func mainPageHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		defer r.Body.Close()
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
		linkStore[hashedURL] = longURL
		shortURL := baseURL + "/" + hashedURL

		rw.WriteHeader(http.StatusCreated)
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err = rw.Write([]byte(shortURL))
		if err != nil {
			return
		}
	} else if r.Method == http.MethodGet {
		shortURL := r.URL.Path[1:]

		if shortURL == "" {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		rw.Header().Set("Location", linkStore[shortURL])
		rw.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

}

func shorten(longURL string) string {
	hashedString := sha1.Sum([]byte(longURL))
	return hex.EncodeToString(hashedString[:4])
}
