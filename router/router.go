package router

import (
	"github.com/ANiWarlock/yandex_go_url_sh.git/internal/app"
	"github.com/go-chi/chi/v5"
)

func NewShortenerRouter(myApp *app.App) chi.Router {
	r := chi.NewRouter()

	r.Post("/", myApp.GetShortURLHandler)
	r.Get("/{shortURL}", myApp.LongURLRedirectHandler)

	return r
}
