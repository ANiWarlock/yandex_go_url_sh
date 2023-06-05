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
	sugar, err := logger.Initialize("info")
	if err != nil {
		log.Fatalf("Cannot init logger: %v", err)
	}

	cfg, err := config.InitConfig()
	if err != nil {
		sugar.Fatalf("Cannot init config: %v", err)
	}
	store, err := storage.InitStorage(*cfg)
	if err != nil {
		sugar.Fatalf("Cannot init storage: %v", err)
	}
	myApp := app.NewApp(cfg, store, sugar)
	shortRouter := router.NewShortenerRouter(myApp, sugar)

	sugar.Infow(
		"Starting server",
		"addr", cfg.Host,
	)
	sugar.Fatal(http.ListenAndServe(cfg.Host, shortRouter))
}
