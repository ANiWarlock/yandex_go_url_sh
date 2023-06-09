package router

import (
	"github.com/ANiWarlock/yandex_go_url_sh.git/internal/app"
	"github.com/ANiWarlock/yandex_go_url_sh.git/router/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewShortenerRouter(myApp *app.App, sugar *zap.SugaredLogger) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.SugarLogger(sugar))
	r.Use(middleware.Gzip)
	r.Use(middleware.SetAuthCookie)

	r.Post("/", myApp.GetShortURLHandler)
	r.Post("/api/shorten", myApp.GetShortURLHandler)
	r.Post("/api/shorten/batch", myApp.APIBatchHandler)
	r.Get("/api/user/urls", myApp.APIUserUrlsHandler)
	r.Delete("/api/user/urls", myApp.APIUserUrlsDeleteHandler)
	r.Get("/{shortURL}", myApp.LongURLRedirectHandler)
	r.Get("/ping", myApp.PingHandler)

	return r
}
