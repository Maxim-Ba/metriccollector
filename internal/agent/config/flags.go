package config

import (
	"flag"
	"os"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
)

type ParsedFlags struct {
	FlagRunAddr        utils.FlagValue[string]
	FlagReportInterval utils.FlagValue[int]
	FlagPollInterval   utils.FlagValue[int]
	// debug info warn error
	LogLevel      utils.FlagValue[string]
	Key           utils.FlagValue[string]
	RateLimit     utils.FlagValue[int]
	CryptoKeyPath utils.FlagValue[string]
	ConfigPath    utils.FlagValue[string]
}

func ParseFlags() *ParsedFlags {
	flags := &ParsedFlags{}

	flag.StringVar(&flags.FlagRunAddr.Value, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flags.FlagReportInterval.Value, "r", 10, "seconds interval to send report")
	flag.IntVar(&flags.FlagPollInterval.Value, "p", 2, "seconds interval to collect metrics")
	flag.StringVar(&flags.LogLevel.Value, "ll", "debug", "log level: debug info warn error")
	flag.StringVar(&flags.Key.Value, "k", "", "private key for signature")
	flag.IntVar(&flags.RateLimit.Value, "l", 10, "simultaneously get metrics")
	flag.StringVar(&flags.CryptoKeyPath.Value, "crypto-key", "", "path for public key for signature")
	flag.StringVar(&flags.ConfigPath.Value, "c", "", "path for configuration by json")

	flag.Parse()

	flag.Usage = func() {
		logger.LogInfo(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	// Check for unknown flags
	for _, arg := range flag.Args() {
		if arg[0] == '-' {
			logger.LogInfo(os.Stderr, "Unknown flag: %s\n", arg)
			flag.Usage()
			os.Exit(1)
		}
	}

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "a":
			flags.FlagRunAddr.Passed = true
		case "r":
			flags.FlagReportInterval.Passed = true
		case "p":
			flags.FlagPollInterval.Passed = true
		case "ll":
			flags.LogLevel.Passed = true
		case "k":
			flags.Key.Passed = true
		case "l":
			flags.RateLimit.Passed = true
		case "crypto-key":
			flags.CryptoKeyPath.Passed = true
		case "c":
			flags.ConfigPath.Passed = true
		}
	})
	return flags
}
