package app

import (
	"github.com/ANiWarlock/yandex_go_url_sh.git/config"
	"github.com/ANiWarlock/yandex_go_url_sh.git/storage"
	"go.uber.org/zap"
)

type App struct {
	cfg     *config.AppConfig
	storage storage.Storage
	sugar   *zap.SugaredLogger
}

func NewApp(cfg *config.AppConfig, store storage.Storage, sugar *zap.SugaredLogger) *App {
	return &App{
		cfg:     cfg,
		storage: store,
		sugar:   sugar,
	}
}
