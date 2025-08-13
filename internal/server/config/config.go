package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/pkg/utils"
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
	TrustedSubnet       string `json:"trusted_subnet"`
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
		Address:             utils.ResolveString(envConfig.Address, flags.RunAddr, fileConfig.Address),
		StoreIntervalSecond: utils.ResolveInt(envConfig.StoreIntervalSecond, flags.StoreIntervalSecond, fileConfig.StoreIntervalSecond),
		StoragePath:         utils.ResolveString(envConfig.StoragePath, flags.StoragePath, fileConfig.StoragePath),
		Restore:             utils.ResolveBool(isRestoreSet(), envConfig.Restore, flags.Restore, fileConfig.Restore),
		LogLevel:            utils.ResolveString(envConfig.LogLevel, flags.LogLevel, fileConfig.LogLevel),
		DatabaseDSN:         utils.ResolveString(envConfig.DatabaseDSN, flags.DatabaseDSN, fileConfig.DatabaseDSN),
		MigrationsPath:      utils.ResolveString(envConfig.MigrationsPath, flags.MigrationsPath, fileConfig.MigrationsPath),
		Key:                 utils.ResolveString(envConfig.Key, flags.Key, fileConfig.Key),
		ProfileFileCPU:      utils.ResolveString(envConfig.ProfileFileCPU, flags.ProfileFileCPU, fileConfig.ProfileFileCPU),
		ProfileFileMem:      utils.ResolveString(envConfig.ProfileFileMem, flags.ProfileFileMem, fileConfig.ProfileFileMem),
		IsProfileOn:         utils.ResolveBool(isProfileOnSet(), envConfig.IsProfileOn, flags.IsProfileOn, fileConfig.IsProfileOn),
		CryptoKeyPath:       utils.ResolveString(envConfig.CryptoKeyPath, flags.CryptoKeyPath, fileConfig.CryptoKeyPath),
		TrustedSubnet:       utils.ResolveString(envConfig.TrustedSubnet, flags.TrustedSubnet, fileConfig.TrustedSubnet),
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
