package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
)

const baseUrl = "http://localhost:8080"

var (
	linkStore map[string]string
)

func mainPage(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		defer r.Body.Close()
		responseData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, fmt.Sprintf("Error in request body: %s", err), http.StatusBadRequest)
			return
		}

		longUrl := string(responseData)
		hashedUrl := shorten(longUrl)
		shortUrl := baseUrl + "/" + hashedUrl
		linkStore[hashedUrl] = longUrl

		rw.Header().Set("content-type", "text/plain")
		rw.WriteHeader(http.StatusCreated)
		_, err = rw.Write([]byte(shortUrl))
		if err != nil {
			return
		}
	} else if r.Method == http.MethodGet {
		shortUrl := r.URL.Path

		if shortUrl == "" {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		rw.Header().Set("Location", linkStore[shortUrl])
		rw.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

}

func shorten(longUrl string) string {
	hashedString := md5.Sum([]byte(longUrl))
	return string(hashedString[:5])
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
