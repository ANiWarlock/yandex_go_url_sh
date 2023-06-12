package storage

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ANiWarlock/yandex_go_url_sh.git/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"os"
)

type Storage struct {
	store    map[string]string
	filename string
	file     *os.File
	writer   *bufio.Writer
	lastUUID int
	db       *sql.DB
}

type Item struct {
	UUID     int    `json:"uuid"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"original_url"`
}

func InitStorage(cfg config.AppConfig) (*Storage, error) {
	storage := Storage{
		store:    make(map[string]string),
		filename: cfg.Filename,
		lastUUID: 1,
	}

	if cfg.Filename != "" {
		file, err := os.OpenFile(cfg.Filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}

		err = storage.loadFromFile(file)
		if err != nil {
			return nil, err
		}

		storage.writer = bufio.NewWriter(file)

	}

	if cfg.DatabaseDSN != "" {
		db, err := sql.Open("pgx", cfg.DatabaseDSN)
		if err != nil {
			return nil, fmt.Errorf("failed to open db connection: %w", err)
			//return &storage, nil
		}
		storage.db = db
		defer storage.db.Close()
	}

	return &storage, nil
}

func (s *Storage) loadFromFile(file *os.File) error {
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var line Item
		err := json.Unmarshal(scanner.Bytes(), &line)
		if err != nil {
			return fmt.Errorf("failed to unmarshal a line from the file storage document: %w", err)
		}

		s.store[line.ShortURL] = line.LongURL
		s.lastUUID = line.UUID
	}

	return nil
}

func (s *Storage) SaveLongURL(hashedURL, longURL string) error {
	exist := s.store[hashedURL] != ""
	if exist {
		return nil
	}

	s.store[hashedURL] = longURL

	if s.filename == "" {
		return nil
	}

	item := Item{
		UUID:     s.lastUUID + 1,
		ShortURL: hashedURL,
		LongURL:  longURL,
	}
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	if _, err := s.writer.Write(data); err != nil {
		return err
	}

	if err := s.writer.WriteByte('\n'); err != nil {
		return err
	}
	s.lastUUID += 1
	return s.writer.Flush()
}

func (s *Storage) GetLongURL(hashedURL string) (string, bool) {
	longURL := s.store[hashedURL]

	if longURL == "" {
		return "", false
	}

	return longURL, true
}

func (s *Storage) PiingDB() bool {
	if s.db != nil {
		if err := s.db.Ping(); err == nil {
			return true
		}
	}

	return false
}
