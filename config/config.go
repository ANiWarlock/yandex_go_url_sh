package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v8"
)

type AppConfig struct {
	Host        string `env:"SERVER_ADDRESS"`
	BaseURL     string `env:"BASE_URL"`
	Filename    string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN string `env:"DATABASE_DSN"`
}

func InitConfig() (*AppConfig, error) {
	cfg := AppConfig{}
	cfg.parseFlags()
	err := cfg.parseEnvs()
	if err != nil {
		return &cfg, err
	}
	return &cfg, nil
}

func (c *AppConfig) parseFlags() {
	flag.StringVar(&c.Host, "a", "localhost:8080", "Укажите адрес в формате host:port")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "Укажите возвращаемый адрес в формате http://host:port")
	flag.StringVar(&c.Filename, "f", "/tmp/short-url-db.json", "Полное имя файла, куда сохраняются данные в формате JSON")
	flag.StringVar(&c.DatabaseDSN, "d", "host=localhost port=5433 user=postgres password=postgres dbname=postgres sslmode=disable", "Параметры подключения к БД")
	flag.Parse()
}

func (c *AppConfig) parseEnvs() error {
	if err := env.Parse(c); err != nil {
		return fmt.Errorf("failed to parse env vars: %w", err)
	}
	return nil
}
