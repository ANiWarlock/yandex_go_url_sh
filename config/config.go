package config

import (
	"flag"
	"github.com/caarlos0/env/v8"
)

type AppConfig struct {
	Host    string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
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
	flag.Parse()
}

func (c *AppConfig) parseEnvs() error {
	if err := env.Parse(&c); err != nil {
		return err
	}
	return nil
}
