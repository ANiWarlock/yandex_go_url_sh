package main

import (
	"github.com/ANiWarlock/yandex_go_url_sh.git/internal/app"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func shortenerRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/", app.MainPageHandler)
	r.Get("/{shortURL}", app.LongUrlRedirectHandler)

	return r
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", shortenerRouter()))
}
