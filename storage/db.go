package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ANiWarlock/yandex_go_url_sh.git/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	db *sql.DB
}

func InitDBStorage(cfg config.AppConfig) (*DBStorage, error) {
	dbStore := DBStorage{}

	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return &dbStore, fmt.Errorf("failed to open db connection: %w", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS urls (id SERIAL PRIMARY KEY , short varchar(16) NOT NULL, long varchar(255) NOT NULL)")
	if err != nil {
		return &dbStore, fmt.Errorf("failed to create db table: %w", err)
	}

	_, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS  long_url_idx ON urls (long)")
	if err != nil {
		return &dbStore, fmt.Errorf("failed to create unique index: %w", err)
	}

	dbStore.db = db

	return &dbStore, nil
}

func (dbs *DBStorage) SaveLongURL(hashedURL, longURL string) error {
	item := Item{
		ShortURL: hashedURL,
		LongURL:  longURL,
	}

	result, err := dbs.db.Exec("INSERT INTO urls (short, long) VALUES ($1, $2) ON CONFLICT (long) DO NOTHING;", item.ShortURL, item.LongURL)
	if err != nil {
		return fmt.Errorf("cannot save to db: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot save to db: %w", err)
	}
	if count == 0 {
		return &UniqueViolationError{ShortURL: item.ShortURL}
	}
	return nil
}

func (dbs *DBStorage) GetLongURL(hashedURL string) (*Item, error) {
	var itemFromDB Item
	row := dbs.db.QueryRow("SELECT id, short, long  FROM urls WHERE short = $1;", hashedURL)
	err := row.Scan(&itemFromDB.UUID, &itemFromDB.ShortURL, &itemFromDB.LongURL)
	if err == sql.ErrNoRows {
		return nil, errors.New("longURL not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed scanning row: %w", err)
	}

	return &itemFromDB, nil
}

func (dbs *DBStorage) Ping() error {
	if err := dbs.db.Ping(); err != nil {
		return err
	}

	return nil
}

func (dbs *DBStorage) CloseDB() error {
	if dbs.db == nil {
		return nil
	}
	if err := dbs.db.Close(); err != nil {
		return fmt.Errorf("cannot close db: %w", err)
	}
	return nil
}

type UniqueViolationError struct {
	ShortURL string
}

func (uve *UniqueViolationError) Error() string {
	return uve.ShortURL
}
