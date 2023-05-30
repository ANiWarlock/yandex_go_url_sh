package storage

import (
	"bufio"
	"encoding/json"
	"github.com/ANiWarlock/yandex_go_url_sh.git/config"
	"github.com/ANiWarlock/yandex_go_url_sh.git/logger"
	"os"
)

type Storage struct {
	store    map[string]string
	filename string
	file     *os.File
	writer   *bufio.Writer
	lastUUID int
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
			return nil, err
		}

		err = storage.loadFromFile(file)
		if err != nil {
			return nil, err
		}

		storage.writer = bufio.NewWriter(file)
	} else {
		logger.Sugar.Infoln("Filename не указан, пропускаем загрузку из файла")
	}

	return &storage, nil
}

func (s *Storage) loadFromFile(file *os.File) error {
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var line Item
		err := json.Unmarshal(scanner.Bytes(), &line)
		if err != nil {
			return err
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

	if s.filename != "" {
		item := Item{
			UUID:     s.lastUUID + 1,
			ShortURL: hashedURL,
			LongURL:  longURL,
		}
		data, err := json.Marshal(item)
		if err != nil {
			return err
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

	return nil
}

func (s *Storage) GetLongURL(hashedURL string) (string, bool) {
	longURL := s.store[hashedURL]

	if longURL == "" {
		return "", false
	}

	return longURL, true
}
