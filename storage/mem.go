package storage

import (
	"errors"
	"fmt"
	"github.com/ANiWarlock/yandex_go_url_sh.git/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type MemStorage struct {
	store map[string]string
}

func InitMemStorage(cfg config.AppConfig) *MemStorage {
	memStore := MemStorage{
		store: make(map[string]string),
	}

	return &memStore
}

func (ms *MemStorage) SaveLongURL(hashedURL, longURL string) error {
	if ms.store[hashedURL] != "" {
		return nil
	}

	ms.store[hashedURL] = longURL
	return nil
}

func (ms *MemStorage) BatchInsert(items []Item) error {
	for _, item := range items {
		err := ms.SaveLongURL(item.ShortURL, item.LongURL)
		if err != nil {
			return fmt.Errorf("failed to batch save item: %w", err)
		}
	}
	return nil
}

func (ms *MemStorage) GetLongURL(hashedURL string) (*Item, error) {
	item := Item{
		ShortURL: hashedURL,
	}

	longURL := ms.store[hashedURL]

	if longURL == "" {
		return &item, errors.New("longURL not found")
	}

	item.LongURL = longURL
	return &item, nil
}

func (ms *MemStorage) Ping() error {
	return nil
}

func (ms *MemStorage) CloseDB() error {
	return nil
}
