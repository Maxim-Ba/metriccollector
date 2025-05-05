package config

import (
	"flag"
)

type ParsedFlags struct {
	FlagRunAddr             string
	FlagStoreIntervalSecond int
	FlagStoragePath         string
	FlagRestore             bool
	// debug info warn error
	LogLevel       string
	DatabaseDSN    string
	MigrationsPath string
}

var FlagRunAddr string
var FlagStoreIntervalSecond int
var FlagStoragePath string
var FlagRestore bool
var LogLevel string
var DatabaseDSN string
var MigrationsPath string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() *ParsedFlags {
	// регистрируем переменную FlagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(&FlagStoreIntervalSecond, "i", 300, "interval after save metrics to file")
	flag.StringVar(&FlagStoragePath, "f", "./store.json", "file path for save metrics")
	flag.BoolVar(&FlagRestore, "r", true, "load metrics at server start ")
	flag.StringVar(&LogLevel, "l", "debug", "log level: debug info warn error")
	flag.StringVar(&DatabaseDSN, "d", "", "addres to connect to database")          //-d postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
	flag.StringVar(&MigrationsPath, "m", "migrations", "path to migrations") // -m migrations
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
	return &ParsedFlags{
		FlagRunAddr:             FlagRunAddr,
		FlagStoreIntervalSecond: FlagStoreIntervalSecond,
		FlagStoragePath:         FlagStoragePath,
		FlagRestore:             FlagRestore,
		LogLevel:                LogLevel,
		DatabaseDSN:             DatabaseDSN,
		MigrationsPath:          MigrationsPath,
	}
}
