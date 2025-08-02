package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Parameters struct {
	Address             string `json:"address"`
	StoreIntervalSecond int    `json:"store_interval"`
	StoragePath         string `json:"store_file"`
	Restore             bool   `json:"restore"`
	LogLevel            string `json:"log_level"`
	DatabaseDSN         string `json:"database_dsn"`
	MigrationsPath      string `json:"migrations_path"`
	Key                 string `json:"key"`
	ProfileFileCPU      string `json:"profile_file_cpu"`
	ProfileFileMem      string `json:"profile_file_mem"`
	IsProfileOn         bool   `json:"is_profile_on"`
	CryptoKeyPath       string `json:"crypto_key"`
}

func New() Parameters {
	flags := ParseFlags()
	envConfig := ParseEnv()

	fileConfig := getParamsByConfigPath(resolveString(envConfig.ConfigPath, flags.ConfigPath, ""))

	parameters := Parameters{
		Address:             resolveString(envConfig.Address, flags.RunAddr, fileConfig.Address),
		StoreIntervalSecond: resolveInt(envConfig.StoreIntervalSecond, flags.StoreIntervalSecond, fileConfig.StoreIntervalSecond),
		StoragePath:         resolveString(envConfig.StoragePath, flags.StoragePath, fileConfig.StoragePath),
		Restore:             resolveBool(isRestoreSet(), envConfig.Restore, flags.Restore, fileConfig.Restore),
		LogLevel:            resolveString(envConfig.LogLevel, flags.LogLevel, fileConfig.LogLevel),
		DatabaseDSN:         resolveString(envConfig.DatabaseDSN, flags.DatabaseDSN, fileConfig.DatabaseDSN),
		MigrationsPath:      resolveString(envConfig.MigrationsPath, flags.MigrationsPath, fileConfig.MigrationsPath),
		Key:                 resolveString(envConfig.Key, flags.Key, fileConfig.Key),
		ProfileFileCPU:      resolveString(envConfig.ProfileFileCPU, flags.ProfileFileCPU, fileConfig.ProfileFileCPU),
		ProfileFileMem:      resolveString(envConfig.ProfileFileMem, flags.ProfileFileMem, fileConfig.ProfileFileMem),
		IsProfileOn:         resolveBool(isProfileOnSet(), envConfig.IsProfileOn, flags.IsProfileOn, fileConfig.IsProfileOn),
		CryptoKeyPath:       resolveString(envConfig.CryptoKeyPath, flags.CryptoKeyPath, fileConfig.CryptoKeyPath),
	}
	fmt.Printf("%+v\n", parameters)
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
