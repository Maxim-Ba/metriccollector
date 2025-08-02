package config

import (
	"encoding/json"
	"fmt"
	"os"
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
	fileConfig := getParamsByConfigPath(resolveString(envConfig.ConfigPath, flags.ConfigPath, ""))
	parameters := Parameters{
		Address:        resolveString(envConfig.Address, flags.FlagRunAddr, fileConfig.Address),
		ReportInterval: resolveInt(envConfig.ReportInterval, flags.FlagReportInterval, fileConfig.ReportInterval),
		LogLevel:       resolveString(envConfig.LogLevel, flags.LogLevel, fileConfig.LogLevel),
		PollInterval:   resolveInt(envConfig.PollInterval, flags.FlagPollInterval, fileConfig.PollInterval),
		RateLimit:      resolveInt(envConfig.RateLimit, flags.RateLimit, fileConfig.RateLimit),
		Key:            resolveString(envConfig.Key, flags.Key, fileConfig.Key),

		CryptoKeyPath: resolveString(envConfig.CryptoKeyPath, flags.CryptoKeyPath, fileConfig.CryptoKeyPath),
	}
	return parameters
}

func getParamsByConfigPath(configPath string) Parameters {
	var parameters Parameters
	if configPath == "" {
		return parameters
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("%v", fmt.Errorf("read config path file: %v", err))
		return parameters
	}
	err = json.Unmarshal(data, &parameters)
	if err != nil {
		fmt.Printf("%v", fmt.Errorf("unmarshal params: %v", err))
		return parameters
	}
	return parameters
}

func resolveString(envValue string, flag FlagValue[string], fileValue string) string {
	if envValue != "" {
		return envValue
	}
	if flag.Passed {
		return flag.Value
	}
	if fileValue != "" {
		return fileValue
	}
	return flag.Value
}

func resolveInt(envValue int, flag FlagValue[int], fileValue int) int {
	if envValue != 0 {
		return envValue
	}
	if flag.Passed {
		return flag.Value
	}
	if fileValue != 0 {
		return fileValue
	}
	return flag.Value
}

func resolveBool(isEnvSet bool, envValue bool, flag FlagValue[bool], fileValue bool) bool {
	if isEnvSet {
		return envValue
	}
	if flag.Passed {
		return flag.Value
	}
	if fileValue {
		return fileValue
	}
	return flag.Value
}
