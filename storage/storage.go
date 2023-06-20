package storage

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ANiWarlock/yandex_go_url_sh.git/config"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

	if cfg.DatabaseDSN != "" {
		db, err := sql.Open("pgx", cfg.DatabaseDSN)
		if err != nil {
			return nil, fmt.Errorf("failed to open db connection: %w", err)
		}

		_, err = db.Exec("CREATE TABLE IF NOT EXISTS urls (id SERIAL PRIMARY KEY , short varchar(16) NOT NULL, long varchar(255) NOT NULL)")
		if err != nil {
			return nil, fmt.Errorf("failed to create db table: %w", err)
		}

		_, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS  long_url_idx ON urls (long)")
		if err != nil {
			return nil, fmt.Errorf("failed to create unique index: %w", err)
		}

		storage.db = db

		return &storage, nil
	}

	if cfg.Filename == "" {
		return &storage, nil
	}

	file, err := os.OpenFile(cfg.Filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	err = storage.loadFromFile(file)
	if err != nil {
		return nil, err
	}

	storage.writer = bufio.NewWriter(file)

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
	item := Item{
		ShortURL: hashedURL,
		LongURL:  longURL,
	}

	if s.db != nil {
		if err := s.saveItemToDB(item); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr); pgErr.Code == pgerrcode.UniqueViolation {
				getItem, getErr := s.getItemFromDB(longURL, "long")
				if getErr != nil {
					return getErr
				}
				return &UniqueViolationError{Err: err, ShortURL: getItem.ShortURL}
			}
			return err
		}

		return nil
	}

	exist := s.store[hashedURL] != ""
	if exist {
		return nil
	}

	s.store[hashedURL] = longURL

	if s.filename == "" {
		return nil
	}

	item.UUID = s.lastUUID + 1
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

func (s *Storage) GetLongURL(hashedURL string) (string, error) {
	if s.db != nil {
		itemFromDB, err := s.getItemFromDB(hashedURL, "short")
		if err != nil {
			return "", err
		}

		if itemFromDB == (Item{}) {
			return "", errors.New("longURL not found")
		}

		return itemFromDB.LongURL, nil
	}

	longURL := s.store[hashedURL]

	if longURL == "" {
		return "", errors.New("longURL not found")
	}

	return longURL, nil
}

func (s *Storage) PingDB() error {
	if s.db == nil {
		return errors.New("DB not configured")
	}
	if err := s.db.Ping(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) getItemFromDB(longURL, field string) (Item, error) {
	var item Item
	row := s.db.QueryRow("SELECT id, short, long  FROM urls WHERE "+field+" = $1;", longURL)
	err := row.Scan(&item.UUID, &item.ShortURL, &item.LongURL)
	if err == sql.ErrNoRows {
		return item, nil
	}
	if err != nil {
		return item, fmt.Errorf("failed scanning row: %w", err)
	}

	return item, nil
}

func (s *Storage) saveItemToDB(item Item) error {
	_, err := s.db.Exec("INSERT INTO urls (short, long) VALUES ($1, $2);", item.ShortURL, item.LongURL)
	if err != nil {
		return fmt.Errorf("cannot save to db: %w", err)
	}
	return nil
}

func (s *Storage) CloseDB() error {
	if s.db == nil {
		return nil
	}
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("cannot close db: %w", err)
	}
	return nil
}

type UniqueViolationError struct {
	Err      error
	ShortURL string
}

func (uve *UniqueViolationError) Error() string {
	return fmt.Sprintf("%s | %v", uve.ShortURL, uve.Err)
}
