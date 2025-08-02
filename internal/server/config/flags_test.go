package config

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func TestParseFlags(t *testing.T) {
	// Сохраняем оригинальные аргументы командной строки и флаги
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tests := []struct {
		name     string
		args     []string
		expected ParsedFlags
	}{
		{
			name: "default values",
			args: []string{"test"},
			expected: ParsedFlags{
				RunAddr:             FlagValue[string]{Value: ":8080"},
				StoreIntervalSecond: FlagValue[int]{Value: 300},
				StoragePath:         FlagValue[string]{Value: "./store.json"},
				Restore:             FlagValue[bool]{Value: true},
				LogLevel:            FlagValue[string]{Value: "debug"},
				DatabaseDSN:         FlagValue[string]{Value: ""},
				MigrationsPath:      FlagValue[string]{Value: "migrations"},
				Key:                 FlagValue[string]{Value: ""},
				ProfileFileCPU:      FlagValue[string]{Value: ""},
				ProfileFileMem:      FlagValue[string]{Value: ""},
				IsProfileOn:         FlagValue[bool]{Value: false},
				CryptoKeyPath:       FlagValue[string]{Value: ""},
				ConfigPath:          FlagValue[string]{Value: ""},
			},
		},
		{
			name: "custom values",
			args: []string{
				"test",
				"-a", ":9090",
				"-i", "60",
				"-f", "/tmp/store.json",
				"-r=false",
				"-l", "info",
				"-d", "postgres://user:pass@localhost:5432/db",
				"-m", "custom_migrations",
				"-k", "secret_key",
				"-cpu", "cpu.prof",
				"-mem", "mem.prof",
				"-p=true",
				"-crypto-key", "/path/to/key",
				"-c", "/path/to/config",
			},
			expected: ParsedFlags{
				RunAddr:             FlagValue[string]{Passed: true, Value: ":9090"},
				StoreIntervalSecond: FlagValue[int]{Passed: true, Value: 60},
				StoragePath:         FlagValue[string]{Passed: true, Value: "/tmp/store.json"},
				Restore:             FlagValue[bool]{Passed: true, Value: false},
				LogLevel:            FlagValue[string]{Passed: true, Value: "info"},
				DatabaseDSN:         FlagValue[string]{Passed: true, Value: "postgres://user:pass@localhost:5432/db"},
				MigrationsPath:      FlagValue[string]{Passed: true, Value: "custom_migrations"},
				Key:                 FlagValue[string]{Passed: true, Value: "secret_key"},
				ProfileFileCPU:      FlagValue[string]{Passed: true, Value: "cpu.prof"},
				ProfileFileMem:      FlagValue[string]{Passed: true, Value: "mem.prof"},
				IsProfileOn:         FlagValue[bool]{Passed: true, Value: true},
				CryptoKeyPath:       FlagValue[string]{Passed: true, Value: "/path/to/key"},
				ConfigPath:          FlagValue[string]{Passed: true, Value: "/path/to/config"},
			},
		},
		{
			name: "partial custom values",
			args: []string{
				"test",
				"-a", ":9090",
				"-i", "60",
				"-r=false",
			},
			expected: ParsedFlags{
				RunAddr:             FlagValue[string]{Passed: true, Value: ":9090"},
				StoreIntervalSecond: FlagValue[int]{Passed: true, Value: 60},
				StoragePath:         FlagValue[string]{Value: "./store.json"},
				Restore:             FlagValue[bool]{Passed: true, Value: false},
				LogLevel:            FlagValue[string]{Value: "debug"},
				DatabaseDSN:         FlagValue[string]{Value: ""},
				MigrationsPath:      FlagValue[string]{Value: "migrations"},
				Key:                 FlagValue[string]{Value: ""},
				ProfileFileCPU:      FlagValue[string]{Value: ""},
				ProfileFileMem:      FlagValue[string]{Value: ""},
				IsProfileOn:         FlagValue[bool]{Value: false},
				CryptoKeyPath:       FlagValue[string]{Value: ""},
				ConfigPath:          FlagValue[string]{Value: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сбрасываем флаги перед каждым тестом
			flag.CommandLine = flag.NewFlagSet(tt.args[0], flag.ContinueOnError)
			os.Args = tt.args

			// Вызываем тестируемую функцию
			got := ParseFlags()

			// Проверяем результаты
			if !reflect.DeepEqual(*got, tt.expected) {
				t.Errorf("ParseFlags() = %+v, want %+v", *got, tt.expected)
			}
		})
	}
}

func Test_isFlagPassed(t *testing.T) {
	// Сохраняем оригинальные аргументы командной строки и флаги
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tests := []struct {
		name     string
		args     []string
		flagName string
		want     bool
	}{
		{
			name:     "flag passed",
			args:     []string{"test", "-a", ":9090"},
			flagName: "a",
			want:     true,
		},
		{
			name:     "flag not passed",
			args:     []string{"test"},
			flagName: "a",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сбрасываем флаги перед каждым тестом
			flag.CommandLine = flag.NewFlagSet(tt.args[0], flag.ContinueOnError)
			os.Args = tt.args

			// Регистрируем флаги, как в ParseFlags()
			var runAddr string
			flag.StringVar(&runAddr, "a", ":8080", "address and port to run server")

			// Парсим флаги
			flag.Parse()

			// Проверяем, был ли флаг передан
			got := isFlagPassed(tt.flagName)

			// Проверяем результаты
			if got != tt.want {
				t.Errorf("isFlagPassed() = %v, want %v", got, tt.want)
			}
		})
	}
}
