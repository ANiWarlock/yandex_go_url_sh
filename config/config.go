package config

import (
	"flag"
)

var urls struct {
	host       string
	returnHost string
}

func ParseFlags() {
	flag.StringVar(&urls.host, "a", "localhost:8080", "Укажите адрес в формате host:port")
	flag.StringVar(&urls.returnHost, "b", "http://localhost:8080", "Укажите возвращаемый адрес в формате http://host:port")
	flag.Parse()
}

func GetHost() string {
	return urls.host
}

func GetReturnHost() string {
	return urls.returnHost
}
