package config

type Parameters struct {
	Address             string
	StoreIntervalSecond int
	StoragePath         string
	Restore             bool
	LogLevel            string
	DatabaseDSN         string
	MigrationsPath      string
}

func New() Parameters {
	flags := ParseFlags()
	envConfig := ParseEnv()
	address := envConfig.Addres
	storeInterval := envConfig.StoreIntervalSecond
	storagePath := envConfig.StoragePath
	restore := envConfig.Restore
	logLevel := envConfig.LogLevel
	databaseDSN := envConfig.DatabaseDSN
	migrationsPath := envConfig.MigrationsPath
	if address == "" {
		address = flags.FlagRunAddr
	}
	if !isIntervalSet() {
		storeInterval = flags.FlagStoreIntervalSecond
	}
	if storagePath == "" {
		storagePath = flags.FlagStoragePath
	}
	if !isRestoreSet() {
		restore = flags.FlagRestore
	}
	if logLevel == "" {
		logLevel = flags.LogLevel
	}
	if databaseDSN == "" {
		databaseDSN = flags.DatabaseDSN
	}
	if !isMigrationsPathSet() {
		migrationsPath = flags.MigrationsPath
	}
	return Parameters{
		Address:             address,
		StoreIntervalSecond: storeInterval,
		StoragePath:         storagePath,
		Restore:             restore,
		LogLevel:            logLevel,
		DatabaseDSN:         databaseDSN,
		MigrationsPath:      migrationsPath,
	}
}
