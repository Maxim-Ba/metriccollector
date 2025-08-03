package utils

type FlagValue[T any] struct {
	Passed bool // был ли флаг передан
	Value  T    // текущее значение
}

func ResolveString(envValue string, flag FlagValue[string], fileValue string) string {
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

func ResolveInt(envValue int, flag FlagValue[int], fileValue int) int {
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

func ResolveBool(isEnvSet bool, envValue bool, flag FlagValue[bool], fileValue bool) bool {
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
