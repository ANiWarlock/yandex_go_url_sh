package storage

import (
	"context"
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

func (ms *MemStorage) SaveLongURL(ctx context.Context, hashedURL, longURL, userID string) error {
	if ms.store[hashedURL] != "" {
		return nil
	}

	ms.store[hashedURL] = longURL
	return nil
}

func (ms *MemStorage) BatchInsert(ctx context.Context, items []Item) error {
	for _, item := range items {
		err := ms.SaveLongURL(ctx, item.ShortURL, item.LongURL, item.UserID)
		if err != nil {
			return fmt.Errorf("failed to batch save item: %w", err)
		}
	}
	return nil
}

func (ms *MemStorage) GetLongURL(ctx context.Context, hashedURL string) (*Item, error) {
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

func (ms *MemStorage) GetUserItems(ctx context.Context, userID string) ([]Item, error) {
	//реализовано для DB
	return nil, nil
}

func (ms *MemStorage) BatchDeleteURL(ctx context.Context, items []Item) {
	//реализовано для DB
}

func (ms *MemStorage) Ping(ctx context.Context) error {
	return nil
}

func (ms *MemStorage) CloseDB() error {
	return nil
}
