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

	_, err = pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS urls (id SERIAL PRIMARY KEY , short varchar(16) NOT NULL, long varchar(255) NOT NULL UNIQUE, user_id varchar(36))")
	if err != nil {
		return &dbStore, fmt.Errorf("failed to create db table: %w", err)
	}

	dbStore.pool = pool

	return &dbStore, nil
}

func (dbs *DBStorage) SaveLongURL(hashedURL, longURL, userID string) error {
	item := Item{
		ShortURL: hashedURL,
		LongURL:  longURL,
		UserID:   userID,
	}

	result, err := dbs.pool.Exec(context.Background(), "INSERT INTO urls (short, long, user_id) VALUES ($1, $2, $3) ON CONFLICT (long) DO NOTHING;", item.ShortURL, item.LongURL, item.UserID)
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
		batch.Queue("INSERT INTO urls (short, long, user_id) VALUES ($1, $2, $3) ON CONFLICT (long) DO NOTHING;", item.ShortURL, item.LongURL, item.UserID)
	}
	err := dbs.pool.SendBatch(context.Background(), batch).Close()
	if err != nil {
		return fmt.Errorf("SendBatch error: %v", err)
	}

	return nil
}

func (dbs *DBStorage) GetLongURL(hashedURL string) (*Item, error) {
	var itemFromDB Item
	row := dbs.pool.QueryRow(context.Background(), "SELECT id, short, long, user_id  FROM urls WHERE short = $1;", hashedURL)
	err := row.Scan(&itemFromDB.UUID, &itemFromDB.ShortURL, &itemFromDB.LongURL, &itemFromDB.UserID)
	if err == sql.ErrNoRows {
		return nil, errors.New("longURL not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed scanning row: %w", err)
	}

	return &itemFromDB, nil
}

func (dbs *DBStorage) GetUserItems(userID string) ([]Item, error) {
	itemsFromDB := make([]Item, 0)
	rows, err := dbs.pool.Query(context.Background(), "SELECT id, short, long, user_id  FROM urls WHERE user_id = $1;", userID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i Item
		err = rows.Scan(&i.UUID, &i.ShortURL, &i.LongURL, &i.UserID)

		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		itemsFromDB = append(itemsFromDB, i)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("db reading error: %w", err)
	}
	return itemsFromDB, nil
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
