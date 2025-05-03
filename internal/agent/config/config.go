package config

type Parameters struct {
	Addres         string
	ReportInterval int
	PollInterval   int
	LogLevel       string
}

func New() Parameters {
	flags := ParseFlags()
	envConfig := ParseEnv()
	address := envConfig.Address
	pollInterval := envConfig.PollInterval
	reportInterval := envConfig.ReportInterval
	logLevel := envConfig.LogLevel

	if address == "" {
		address = flags.FlagRunAddr
	}
	if pollInterval == 0 {
		pollInterval = flags.FlagPollInterval
	}
	if reportInterval == 0 {
		reportInterval = flags.FlagReportInterval
	}
	if logLevel == "" {
		logLevel = flags.LogLevel
	}
	return Parameters{
		Addres:         address,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
		LogLevel:       logLevel,
	}
}
