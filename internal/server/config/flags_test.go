package config

import (
	"flag"
	"os"
	"reflect"
	"testing"

	"github.com/Maxim-Ba/metriccollector/pkg/utils"
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
				RunAddr:             utils.FlagValue[string]{Value: ":8080"},
				StoreIntervalSecond: utils.FlagValue[int]{Value: 300},
				StoragePath:         utils.FlagValue[string]{Value: "./store.json"},
				Restore:             utils.FlagValue[bool]{Value: true},
				LogLevel:            utils.FlagValue[string]{Value: "debug"},
				DatabaseDSN:         utils.FlagValue[string]{Value: ""},
				MigrationsPath:      utils.FlagValue[string]{Value: "migrations"},
				Key:                 utils.FlagValue[string]{Value: ""},
				ProfileFileCPU:      utils.FlagValue[string]{Value: ""},
				ProfileFileMem:      utils.FlagValue[string]{Value: ""},
				IsProfileOn:         utils.FlagValue[bool]{Value: false},
				CryptoKeyPath:       utils.FlagValue[string]{Value: ""},
				ConfigPath:          utils.FlagValue[string]{Value: ""},
				GrpcServer:          utils.FlagValue[string]{Value: ":8081"},
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
				"-ga", ":9001",
			},
			expected: ParsedFlags{
				RunAddr:             utils.FlagValue[string]{Passed: true, Value: ":9090"},
				StoreIntervalSecond: utils.FlagValue[int]{Passed: true, Value: 60},
				StoragePath:         utils.FlagValue[string]{Passed: true, Value: "/tmp/store.json"},
				Restore:             utils.FlagValue[bool]{Passed: true, Value: false},
				LogLevel:            utils.FlagValue[string]{Passed: true, Value: "info"},
				DatabaseDSN:         utils.FlagValue[string]{Passed: true, Value: "postgres://user:pass@localhost:5432/db"},
				MigrationsPath:      utils.FlagValue[string]{Passed: true, Value: "custom_migrations"},
				Key:                 utils.FlagValue[string]{Passed: true, Value: "secret_key"},
				ProfileFileCPU:      utils.FlagValue[string]{Passed: true, Value: "cpu.prof"},
				ProfileFileMem:      utils.FlagValue[string]{Passed: true, Value: "mem.prof"},
				IsProfileOn:         utils.FlagValue[bool]{Passed: true, Value: true},
				CryptoKeyPath:       utils.FlagValue[string]{Passed: true, Value: "/path/to/key"},
				ConfigPath:          utils.FlagValue[string]{Passed: true, Value: "/path/to/config"},
				GrpcServer:          utils.FlagValue[string]{Passed: true, Value: ":9001"},
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
				RunAddr:             utils.FlagValue[string]{Passed: true, Value: ":9090"},
				StoreIntervalSecond: utils.FlagValue[int]{Passed: true, Value: 60},
				StoragePath:         utils.FlagValue[string]{Value: "./store.json"},
				Restore:             utils.FlagValue[bool]{Passed: true, Value: false},
				LogLevel:            utils.FlagValue[string]{Value: "debug"},
				DatabaseDSN:         utils.FlagValue[string]{Value: ""},
				MigrationsPath:      utils.FlagValue[string]{Value: "migrations"},
				Key:                 utils.FlagValue[string]{Value: ""},
				ProfileFileCPU:      utils.FlagValue[string]{Value: ""},
				ProfileFileMem:      utils.FlagValue[string]{Value: ""},
				IsProfileOn:         utils.FlagValue[bool]{Value: false},
				CryptoKeyPath:       utils.FlagValue[string]{Value: ""},
				ConfigPath:          utils.FlagValue[string]{Value: ""},
				GrpcServer:          utils.FlagValue[string]{Passed: false, Value: ":8081"},
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
