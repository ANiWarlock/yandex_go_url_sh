package storage

import (
	"fmt"
	"github.com/ANiWarlock/yandex_go_url_sh.git/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage interface {
	SaveLongURL(string, string, string) error
	GetLongURL(string) (*Item, error)
	BatchInsert([]Item) error
	GetUserItems(string) ([]Item, error)
	Ping() error
	CloseDB() error
}

type Item struct {
	UUID     int    `json:"-"`
	UserID   string `json:"-"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"original_url"`
}

var defaultStorage Storage

func InitStorage(cfg config.AppConfig) (Storage, error) {
	var err error
	if cfg.DatabaseDSN != "" {
		defaultStorage, err = InitDBStorage(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to init db storage: %w", err)
		}

		return defaultStorage, nil
	}

	if cfg.Filename != "" {
		defaultStorage, err = InitFileStorage(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to init file storage: %w", err)
		}
		return defaultStorage, nil
	}

	defaultStorage = InitMemStorage(cfg)

	return defaultStorage, nil
}
