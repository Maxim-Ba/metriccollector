package config

type Parameters struct {
	Addres         string
	ReportInterval int
	PollInterval   int
}

func GetParameters() Parameters {
	flags := ParseFlags()
	envConfig := ParseEnv()
	address := envConfig.Address
	pollInterval := envConfig.PollInterval
	reportInterval := envConfig.ReportInterval
	if address == "" {
		address = flags.FlagRunAddr
	}
	if pollInterval == 0 {
		pollInterval = flags.FlagPollInterval
	}
	if reportInterval == 0 {
		reportInterval = flags.FlagReportInterval
	}
	return Parameters{
		Addres:         address,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
	}
}
