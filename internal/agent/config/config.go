package config

type Parameters struct {
	Addres         string
	ReportInterval int
	PollInterval   int
	LogLevel       string
	Key            string
	RateLimit      int
}

func New() Parameters {
	flags := ParseFlags()
	envConfig := ParseEnv()
	address := envConfig.Address
	pollInterval := envConfig.PollInterval
	reportInterval := envConfig.ReportInterval
	logLevel := envConfig.LogLevel
	key := envConfig.Key
	rateLimit := envConfig.RateLimit

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
	if key == "" {
		key = flags.Key
	}
	if rateLimit == 0 {
		rateLimit = flags.RateLimit
	}
	return Parameters{
		Addres:         address,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
		LogLevel:       logLevel,
		Key:            key,
		RateLimit:      rateLimit,
	}
}
