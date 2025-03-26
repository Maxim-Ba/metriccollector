package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address string `env:"ADDRESS"`
	ReportInterval int `env:"REPORT_INTERVAL"`
	PollInterval int `env:"POLL_INTERVAL"`
}

func ParseEnv() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
