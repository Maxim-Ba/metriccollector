package config

import (
	"log"
	"os"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Addres              string `env:"ADDRESS"`
	StoreIntervalSecond int    `env:"STORE_INTERVAL"`
	StoragePath         string `env:"FILE_STORAGE_PATH"`
	Restore             bool   `env:"RESTORE"`
	LogLevel            string `env:"LOG_LEVEL"`
	DatabaseDSN         string `env:"DATABASE_DSN"`
	MigrationsPath      string `env:"MIGRATIONS_PATH"`
	Key                 string `env:"KEY"`
	ProfileFileCPU      string `env:"CPU_FILE"`
	ProfileFileMem      string `env:"MEM_FILE"`
	IsProfileOn         bool   `env:"IS_PROFILE_ON"`
}

func ParseEnv() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &cfg
}

func isRestoreSet() bool {
	_, isSet := os.LookupEnv("RESTORE")
	return isSet
}

func isIntervalSet() bool {
	_, isSet := os.LookupEnv("STORE_INTERVAL")
	return isSet
}

func isMigrationsPathSet() bool {
	_, isSet := os.LookupEnv("MIGRATIONS_PATH")
	return isSet
}
func isProfileOnSet() bool {
	_, isSet := os.LookupEnv("IS_PROFILE_ON")
	return isSet
}
