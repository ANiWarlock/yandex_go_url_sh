package main

import (
	"github.com/ANiWarlock/yandex_go_url_sh.git/config"
	"github.com/ANiWarlock/yandex_go_url_sh.git/internal/app"
	"github.com/ANiWarlock/yandex_go_url_sh.git/logger"
	"github.com/ANiWarlock/yandex_go_url_sh.git/router"
	"github.com/ANiWarlock/yandex_go_url_sh.git/storage"
	"log"
	"net/http"
)

func main() {
	err := logger.Initialize("info")
	if err != nil {
		log.Fatalf("Cannot init logger: %s", err)
	}

	cfg, err := config.InitConfig()
	if err != nil {
		logger.Sugar.Fatalf("Cannot init config: %s", err)
	}
	store := storage.NewStorage()
	myApp := app.NewApp(cfg, store)
	shortRouter := router.NewShortenerRouter(myApp)

	logger.Sugar.Infow(
		"Starting server",
		"addr", cfg.Host,
	)
	logger.Sugar.Fatal(http.ListenAndServe(cfg.Host, shortRouter))
}
