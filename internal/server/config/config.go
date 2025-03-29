package config

type Parameters struct {
	Addres string
}

func GetParameters() Parameters {
	flags := ParseFlags()
	envConfig := ParseEnv()
	addres := envConfig.Addres
	if addres != "" {
		return Parameters{
			Addres: addres,
		}
	}
	return Parameters{
		Addres: flags.FlagRunAddr,
	}
}
