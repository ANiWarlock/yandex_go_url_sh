package router

import (
	"github.com/ANiWarlock/yandex_go_url_sh.git/internal/app"
	"github.com/ANiWarlock/yandex_go_url_sh.git/router/middleware"
	"github.com/go-chi/chi/v5"
)

func NewShortenerRouter(myApp *app.App) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.SugarLogger)
	r.Use(middleware.Gzip)

	r.Post("/", myApp.GetShortURLHandler)
	r.Post("/api/shorten", myApp.APIGetShortURLHandler)
	r.Get("/{shortURL}", myApp.LongURLRedirectHandler)

	return r
}
