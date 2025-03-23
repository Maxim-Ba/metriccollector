package main

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ADDRESS string `env:"ADDRESS"`
}

func parseEnv() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
