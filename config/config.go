package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v8"
)

var urls struct {
	Host    string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
}

func parseFlags() {
	flag.StringVar(&urls.Host, "a", "localhost:8080", "Укажите адрес в формате host:port")
	flag.StringVar(&urls.BaseURL, "b", "http://localhost:8080", "Укажите возвращаемый адрес в формате http://host:port")
	flag.Parse()
}

func parseEnvs() {
	if err := env.Parse(&urls); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func InitOptions() {
	parseFlags()
	parseEnvs()
}

func GetHost() string {
	return urls.Host
}

func GetBaseURL() string {
	return urls.BaseURL
}
