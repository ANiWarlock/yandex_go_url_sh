package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ANiWarlock/yandex_go_url_sh.git/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStorage struct {
	pool *pgxpool.Pool
}

var ErrUniqueViolation = errors.New("unique violation error")

func InitDBStorage(cfg config.AppConfig) (*DBStorage, error) {
	dbStore := DBStorage{}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		return &dbStore, fmt.Errorf("failed to open db connection: %w", err)
	}

	_, err = pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS urls (id SERIAL PRIMARY KEY , short varchar(16) NOT NULL, long varchar(255) NOT NULL UNIQUE)")
	if err != nil {
		return &dbStore, fmt.Errorf("failed to create db table: %w", err)
	}

	dbStore.pool = pool

	return &dbStore, nil
}

func (dbs *DBStorage) SaveLongURL(hashedURL, longURL string) error {
	item := Item{
		ShortURL: hashedURL,
		LongURL:  longURL,
	}

	result, err := dbs.pool.Exec(context.Background(), "INSERT INTO urls (short, long) VALUES ($1, $2) ON CONFLICT (long) DO NOTHING;", item.ShortURL, item.LongURL)
	if err != nil {
		return fmt.Errorf("cannot save to db: %w", err)
	}

	count := result.RowsAffected()
	if count == 0 {
		return fmt.Errorf("%w: %s", ErrUniqueViolation, item.ShortURL)
	}
	return nil
}

func (dbs *DBStorage) BatchInsert(items []Item) error {
	batch := &pgx.Batch{}

	for _, item := range items {
		batch.Queue("INSERT INTO urls (short, long) VALUES ($1, $2) ON CONFLICT (long) DO NOTHING;", item.ShortURL, item.LongURL)
	}
	err := dbs.pool.SendBatch(context.Background(), batch).Close()
	if err != nil {
		return fmt.Errorf("SendBatch error: %v", err)
	}

	return nil
}

func (dbs *DBStorage) GetLongURL(hashedURL string) (*Item, error) {
	var itemFromDB Item
	row := dbs.pool.QueryRow(context.Background(), "SELECT id, short, long  FROM urls WHERE short = $1;", hashedURL)
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
	if err := dbs.pool.Ping(context.Background()); err != nil {
		return err
	}

	return nil
}

func (dbs *DBStorage) CloseDB() error {
	if dbs.pool == nil {
		return nil
	}
	dbs.pool.Close()
	return nil
}
