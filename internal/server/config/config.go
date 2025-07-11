package config

type Parameters struct {
	Address             string
	StoreIntervalSecond int
	StoragePath         string
	Restore             bool
	LogLevel            string
	DatabaseDSN         string
	MigrationsPath      string
	Key                 string
	ProfileFileCPU      string
	ProfileFileMem      string
	IsProfileOn         bool
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
	profileFileCPU := envConfig.ProfileFileCPU
	profileFileMem := envConfig.ProfileFileMem
	isProfileOn := envConfig.IsProfileOn
	key := envConfig.Key
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
	if key == "" {
		key = flags.Key
	}
	if profileFileCPU == "" {
		profileFileCPU = flags.ProfileFileCPU
	}
	if profileFileMem == "" {
		profileFileMem = flags.ProfileFileMem
	}

	if !isProfileOnSet() {
		isProfileOn = flags.IsProfileOn
	}
	return Parameters{
		Address:             address,
		StoreIntervalSecond: storeInterval,
		StoragePath:         storagePath,
		Restore:             restore,
		LogLevel:            logLevel,
		DatabaseDSN:         databaseDSN,
		MigrationsPath:      migrationsPath,
		Key:                 key,
		ProfileFileCPU:      profileFileCPU,
		ProfileFileMem:      profileFileMem,
		IsProfileOn:         isProfileOn,
	}
}
