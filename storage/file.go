package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ANiWarlock/yandex_go_url_sh.git/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"os"
)

type FileStorage struct {
	memStore *MemStorage
	filename string
	file     *os.File
	writer   *bufio.Writer
	lastUUID int
}

func InitFileStorage(cfg config.AppConfig) (*FileStorage, error) {
	memStore := InitMemStorage(cfg)
	fileStore := FileStorage{
		memStore: memStore,
		filename: cfg.Filename,
		lastUUID: 1,
	}

	file, err := os.OpenFile(cfg.Filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	err = fileStore.loadFromFile(file)
	if err != nil {
		return nil, err
	}

	fileStore.writer = bufio.NewWriter(file)

	return &fileStore, nil
}

func (fs *FileStorage) SaveLongURL(hashedURL, longURL, userID string) error {
	item := Item{
		ShortURL: hashedURL,
		LongURL:  longURL,
		UserID:   userID,
	}

	exist := fs.memStore.store[hashedURL] != ""
	if exist {
		return nil
	}

	fs.memStore.store[hashedURL] = longURL

	item.UUID = fs.lastUUID + 1
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	if _, err := fs.writer.Write(data); err != nil {
		return err
	}

	if err := fs.writer.WriteByte('\n'); err != nil {
		return err
	}
	fs.lastUUID += 1
	fs.writer.Flush()
	return nil
}

func (fs *FileStorage) BatchInsert(items []Item) error {
	for _, item := range items {
		err := fs.SaveLongURL(item.ShortURL, item.LongURL, item.UserID)
		if err != nil {
			return fmt.Errorf("failed to batch save item: %w", err)
		}
	}
	return nil
}

func (fs *FileStorage) GetLongURL(hashedURL string) (*Item, error) {
	item := Item{
		ShortURL: hashedURL,
	}

	longURL := fs.memStore.store[hashedURL]

	if longURL == "" {
		return &item, errors.New("longURL not found")
	}
	item.LongURL = longURL
	return &item, nil
}

func (fs *FileStorage) GetUserItems(userID string) ([]Item, error) {
	//реализовано для DB
	return nil, nil
}

func (fs *FileStorage) Ping() error {
	return nil
}

func (fs *FileStorage) loadFromFile(file *os.File) error {
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var line Item
		err := json.Unmarshal(scanner.Bytes(), &line)
		if err != nil {
			return fmt.Errorf("failed to unmarshal a line from the file storage document: %w", err)
		}

		fs.memStore.store[line.ShortURL] = line.LongURL
		fs.lastUUID = line.UUID
	}

	return nil
}

func (fs *FileStorage) CloseDB() error {
	fs.writer.Flush()
	return nil
}
