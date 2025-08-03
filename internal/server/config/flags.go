package config

import (
	"flag"

	"github.com/Maxim-Ba/metriccollector/pkg/utils"
)

type ParsedFlags struct {
	RunAddr             utils.FlagValue[string]
	StoreIntervalSecond utils.FlagValue[int]
	StoragePath         utils.FlagValue[string]
	Restore             utils.FlagValue[bool]
	// debug info warn error
	LogLevel       utils.FlagValue[string]
	DatabaseDSN    utils.FlagValue[string]
	MigrationsPath utils.FlagValue[string]
	Key            utils.FlagValue[string]
	ProfileFileCPU utils.FlagValue[string]
	ProfileFileMem utils.FlagValue[string]
	IsProfileOn    utils.FlagValue[bool]
	CryptoKeyPath  utils.FlagValue[string]
	ConfigPath     utils.FlagValue[string]
}

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() *ParsedFlags {
	flags := &ParsedFlags{}
	flag.StringVar(&flags.RunAddr.Value, "a", ":8080", "address and port to run server")
	flag.IntVar(&flags.StoreIntervalSecond.Value, "i", 300, "interval after save metrics to file")
	flag.StringVar(&flags.StoragePath.Value, "f", "./store.json", "file path for save metrics")
	flag.BoolVar(&flags.Restore.Value, "r", true, "load metrics at server start")
	flag.StringVar(&flags.LogLevel.Value, "l", "debug", "log level: debug info warn error")
	flag.StringVar(&flags.DatabaseDSN.Value, "d", "", "address to connect to database") //-d postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
	flag.StringVar(&flags.MigrationsPath.Value, "m", "migrations", "path to migrations")
	flag.StringVar(&flags.ProfileFileMem.Value, "mem", "", "memory file profile")
	flag.StringVar(&flags.ProfileFileCPU.Value, "cpu", "", "CPU file profile")
	flag.BoolVar(&flags.IsProfileOn.Value, "p", false, "Is profile switched on")
	flag.StringVar(&flags.Key.Value, "k", "", "private key for signature")
	flag.StringVar(&flags.CryptoKeyPath.Value, "crypto-key", "", "path for public key")
	flag.StringVar(&flags.ConfigPath.Value, "c", "", "path for JSON config")

	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "a":
			flags.RunAddr.Passed = true
		case "i":
			flags.StoreIntervalSecond.Passed = true
		case "f":
			flags.StoragePath.Passed = true
		case "r":
			flags.Restore.Passed = true
		case "l":
			flags.LogLevel.Passed = true
		case "d":
			flags.DatabaseDSN.Passed = true
		case "m":
			flags.MigrationsPath.Passed = true
		case "mem":
			flags.ProfileFileMem.Passed = true
		case "cpu":
			flags.ProfileFileCPU.Passed = true
		case "p":
			flags.IsProfileOn.Passed = true
		case "k":
			flags.Key.Passed = true
		case "crypto-key":
			flags.CryptoKeyPath.Passed = true
		case "c":
			flags.ConfigPath.Passed = true
		}
	})
	return flags
}
