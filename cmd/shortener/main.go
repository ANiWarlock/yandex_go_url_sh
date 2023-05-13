package main

import (
	"net/http"
)

const baseURL = "http://localhost:8080"

var linkStore map[string]string

func main() {
	linkStore = make(map[string]string)

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPageHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
