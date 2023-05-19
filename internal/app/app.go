package app

import (
	"github.com/ANiWarlock/yandex_go_url_sh.git/config"
	"github.com/ANiWarlock/yandex_go_url_sh.git/storage"
)

type App struct {
	cfg     *config.AppConfig
	storage *storage.Storage
}

func NewApp(cfg *config.AppConfig, store *storage.Storage) *App {
	return &App{
		cfg:     cfg,
		storage: store,
	}
}
