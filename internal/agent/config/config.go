package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
)

type Parameters struct {
	Address        string `json:"address"`
	ReportInterval int    `json:"report_interval"`
	PollInterval   int    `json:"poll_interval"`
	LogLevel       string `json:"log_level"`
	Key            string `json:"key"`
	RateLimit      int    `json:"rate_limit"`
	CryptoKeyPath  string `json:"crypto_key"`
}

func New() Parameters {
	flags := ParseFlags()
	envConfig := ParseEnv()
	fileConfig, err := getParamsByConfigPath(utils.ResolveString(envConfig.ConfigPath, flags.ConfigPath, ""))
	if err != nil {
		logger.LogError(err)
		return fileConfig
	}
	parameters := Parameters{
		Address:        utils.ResolveString(envConfig.Address, flags.FlagRunAddr, fileConfig.Address),
		ReportInterval: utils.ResolveInt(envConfig.ReportInterval, flags.FlagReportInterval, fileConfig.ReportInterval),
		LogLevel:       utils.ResolveString(envConfig.LogLevel, flags.LogLevel, fileConfig.LogLevel),
		PollInterval:   utils.ResolveInt(envConfig.PollInterval, flags.FlagPollInterval, fileConfig.PollInterval),
		RateLimit:      utils.ResolveInt(envConfig.RateLimit, flags.RateLimit, fileConfig.RateLimit),
		Key:            utils.ResolveString(envConfig.Key, flags.Key, fileConfig.Key),

		CryptoKeyPath: utils.ResolveString(envConfig.CryptoKeyPath, flags.CryptoKeyPath, fileConfig.CryptoKeyPath),
	}
	return parameters
}

func getParamsByConfigPath(configPath string) (Parameters, error) {
	var parameters Parameters
	if configPath == "" {
		return parameters, nil
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return parameters, fmt.Errorf("read config path file: %w", err)
	}
	err = json.Unmarshal(data, &parameters)
	if err != nil {
		return parameters, fmt.Errorf("unmarshal params: %w", err)
	}
	return parameters, nil
}
