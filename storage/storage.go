package storage

import (
	"context"
	"fmt"
	"github.com/ANiWarlock/yandex_go_url_sh.git/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage interface {
	SaveLongURL(context.Context, string, string, string) error
	GetLongURL(context.Context, string) (*Item, error)
	BatchInsert(context.Context, []Item) error
	BatchDeleteURL(context.Context, []Item)
	GetUserItems(context.Context, string) ([]Item, error)
	Ping(context.Context) error
	CloseDB() error
}

type Item struct {
	UUID     int    `json:"-"`
	UserID   string `json:"-"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"original_url"`
	Deleted  bool   `json:"-"`
}

var defaultStorage Storage

func InitStorage(ctx context.Context, cfg config.AppConfig) (Storage, error) {
	var err error
	if cfg.DatabaseDSN != "" {
		defaultStorage, err = InitDBStorage(ctx, cfg)
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
