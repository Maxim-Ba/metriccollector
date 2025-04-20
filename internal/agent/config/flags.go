package config

import (
	"flag"
	"os"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
)

type ParsedFlags struct {
	FlagRunAddr        string
	FlagReportInterval int
	FlagPollInterval   int
	// debug info warn error
	LogLevel string
}

var flagRunAddr string
var flagReportInterval int
var flagPollInterval int
var LogLevel string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() *ParsedFlags {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagReportInterval, "r", 10, "seconds interval to send report")
	flag.IntVar(&flagPollInterval, "p", 2, "seconds interval to collect metrics")
	flag.StringVar(&LogLevel, "l", "debug", "log level: debug info warn error")

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	flag.Usage = func() {
		logger.LogInfo(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Parse the flags
	flag.Parse()

	// Check for unknown flags
	for _, arg := range flag.Args() {
		if arg[0] == '-' {
			logger.LogInfo(os.Stderr, "Unknown flag: %s\n", arg)
			flag.Usage()
			os.Exit(1)
		}
	}
	return &ParsedFlags{
		FlagRunAddr:        flagRunAddr,
		FlagReportInterval: flagReportInterval,
		FlagPollInterval:   flagPollInterval,
		LogLevel:           LogLevel,
	}
}
