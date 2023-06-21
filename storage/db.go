package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ANiWarlock/yandex_go_url_sh.git/config"
	"github.com/jackc/pgx/v5"
)

type DBStorage struct {
	conn *pgx.Conn
}

func InitDBStorage(cfg config.AppConfig) (*DBStorage, error) {
	dbStore := DBStorage{}

	conn, err := pgx.Connect(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		return &dbStore, fmt.Errorf("failed to open db connection: %w", err)
	}

	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS urls (id SERIAL PRIMARY KEY , short varchar(16) NOT NULL, long varchar(255) NOT NULL)")
	if err != nil {
		return &dbStore, fmt.Errorf("failed to create db table: %w", err)
	}

	_, err = conn.Exec(context.Background(), "CREATE UNIQUE INDEX IF NOT EXISTS  long_url_idx ON urls (long)")
	if err != nil {
		return &dbStore, fmt.Errorf("failed to create unique index: %w", err)
	}

	dbStore.conn = conn

	return &dbStore, nil
}

func (dbs *DBStorage) SaveLongURL(hashedURL, longURL string) error {
	item := Item{
		ShortURL: hashedURL,
		LongURL:  longURL,
	}

	result, err := dbs.conn.Exec(context.Background(), "INSERT INTO urls (short, long) VALUES ($1, $2) ON CONFLICT (long) DO NOTHING;", item.ShortURL, item.LongURL)
	if err != nil {
		return fmt.Errorf("cannot save to db: %w", err)
	}

	count := result.RowsAffected()
	if count == 0 {
		return &UniqueViolationError{ShortURL: item.ShortURL}
	}
	return nil
}

func (dbs *DBStorage) BatchInsert(items []Item) error {
	batch := &pgx.Batch{}

	for _, item := range items {
		batch.Queue("INSERT INTO urls (short, long) VALUES ($1, $2) ON CONFLICT (long) DO NOTHING;", item.ShortURL, item.LongURL)
	}
	err := dbs.conn.SendBatch(context.Background(), batch).Close()
	if err != nil {
		return fmt.Errorf("SendBatch error: %v", err)
	}

	return nil
}

func (dbs *DBStorage) GetLongURL(hashedURL string) (*Item, error) {
	var itemFromDB Item
	row := dbs.conn.QueryRow(context.Background(), "SELECT id, short, long  FROM urls WHERE short = $1;", hashedURL)
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
	if err := dbs.conn.Ping(context.Background()); err != nil {
		return err
	}

	return nil
}

func (dbs *DBStorage) CloseDB() error {
	if dbs.conn == nil {
		return nil
	}
	if err := dbs.conn.Close(context.Background()); err != nil {
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
