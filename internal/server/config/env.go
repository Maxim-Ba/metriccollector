package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Addres string `env:"ADDRESS"`
}

func ParseEnv() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
