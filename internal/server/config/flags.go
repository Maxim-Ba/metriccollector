package config

import (
	"flag"
)

type ParsedFlags struct {
	FlagRunAddr string
}

var FlagRunAddr string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() *ParsedFlags {
	// регистрируем переменную FlagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
	return &ParsedFlags{
		FlagRunAddr: FlagRunAddr,
	}
}
