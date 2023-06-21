package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ANiWarlock/yandex_go_url_sh.git/config"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

func (dbs *DBStorage) SaveLongURL(hashedURL, longURL string) (*Item, error) {
	item := Item{
		ShortURL: hashedURL,
		LongURL:  longURL,
	}

	if err := dbs.saveItemToDB(item); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr); pgErr.Code == pgerrcode.UniqueViolation {
			getItem, getErr := dbs.getItemFromDB(longURL, "long")
			if getErr != nil {
				return &item, getErr
			}
			return &item, &UniqueViolationError{Err: err, ShortURL: getItem.ShortURL}
		}
		return &item, err
	}

	return &item, nil

}

func (dbs *DBStorage) GetLongURL(hashedURL string) (*Item, error) {
	itemFromDB, err := dbs.getItemFromDB(hashedURL, "short")
	if err != nil {
		return &Item{}, err
	}

	if itemFromDB == nil {
		return nil, errors.New("longURL not found")
	}

	return itemFromDB, nil

}

func (dbs *DBStorage) Ping() error {
	if err := dbs.db.Ping(); err != nil {
		return err
	}

	return nil
}

func (dbs *DBStorage) getItemFromDB(longURL, field string) (*Item, error) {
	var item Item
	row := dbs.db.QueryRow("SELECT id, short, long  FROM urls WHERE "+field+" = $1;", longURL)
	err := row.Scan(&item.UUID, &item.ShortURL, &item.LongURL)
	if err == sql.ErrNoRows {
		return &item, nil
	}
	if err != nil {
		return &item, fmt.Errorf("failed scanning row: %w", err)
	}

	return &item, nil
}

func (dbs *DBStorage) saveItemToDB(item Item) error {
	_, err := dbs.db.Exec("INSERT INTO urls (short, long) VALUES ($1, $2);", item.ShortURL, item.LongURL)
	if err != nil {
		return fmt.Errorf("cannot save to db: %w", err)
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
	Err      error
	ShortURL string
}

func (uve *UniqueViolationError) Error() string {
	return fmt.Sprintf("%s | %v", uve.ShortURL, uve.Err)
}
