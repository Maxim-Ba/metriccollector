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
				FlagRunAddr:             ":8080",
				FlagStoreIntervalSecond: 300,
				FlagStoragePath:         "./store.json",
				FlagRestore:             true,
				LogLevel:                "debug",
				DatabaseDSN:             "",
				MigrationsPath:          "migrations",
				Key:                     "",
				ProfileFileCPU:          "",
				ProfileFileMem:          "",
				IsProfileOn:             false,
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
			},
			expected: ParsedFlags{
				FlagRunAddr:             ":9090",
				FlagStoreIntervalSecond: 60,
				FlagStoragePath:         "/tmp/store.json",
				FlagRestore:             false,
				LogLevel:                "info",
				DatabaseDSN:             "postgres://user:pass@localhost:5432/db",
				MigrationsPath:          "custom_migrations",
				Key:                     "secret_key",
				ProfileFileCPU:          "cpu.prof",
				ProfileFileMem:          "mem.prof",
				IsProfileOn:             true,
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
				FlagRunAddr:             ":9090",
				FlagStoreIntervalSecond: 60,
				FlagStoragePath:         "./store.json",
				FlagRestore:             false,
				LogLevel:                "debug",
				DatabaseDSN:             "",
				MigrationsPath:          "migrations",
				Key:                     "",
				ProfileFileCPU:          "",
				ProfileFileMem:          "",
				IsProfileOn:             false,
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

func TestParseFlags_GlobalVars(t *testing.T) {
	// Сохраняем оригинальные аргументы командной строки и флаги
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Устанавливаем тестовые аргументы
	os.Args = []string{
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
	}

	// Сбрасываем флаги перед тестом
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Вызываем тестируемую функцию
	_ = ParseFlags()

	// Проверяем, что глобальные переменные установлены правильно
	if FlagRunAddr != ":9090" {
		t.Errorf("FlagRunAddr = %v, want %v", FlagRunAddr, ":9090")
	}
	if FlagStoreIntervalSecond != 60 {
		t.Errorf("FlagStoreIntervalSecond = %v, want %v", FlagStoreIntervalSecond, 60)
	}
	if FlagStoragePath != "/tmp/store.json" {
		t.Errorf("FlagStoragePath = %v, want %v", FlagStoragePath, "/tmp/store.json")
	}
	if FlagRestore != false {
		t.Errorf("FlagRestore = %v, want %v", FlagRestore, false)
	}
	if LogLevel != "info" {
		t.Errorf("LogLevel = %v, want %v", LogLevel, "info")
	}
	if DatabaseDSN != "postgres://user:pass@localhost:5432/db" {
		t.Errorf("DatabaseDSN = %v, want %v", DatabaseDSN, "postgres://user:pass@localhost:5432/db")
	}
	if MigrationsPath != "custom_migrations" {
		t.Errorf("MigrationsPath = %v, want %v", MigrationsPath, "custom_migrations")
	}
	if Key != "secret_key" {
		t.Errorf("Key = %v, want %v", Key, "secret_key")
	}
	if ProfileFileCPU != "cpu.prof" {
		t.Errorf("ProfileFileCPU = %v, want %v", ProfileFileCPU, "cpu.prof")
	}
	if ProfileFileMem != "mem.prof" {
		t.Errorf("ProfileFileMem = %v, want %v", ProfileFileMem, "mem.prof")
	}
	if IsProfileOn != true {
		t.Errorf("IsProfileOn = %v, want %v", IsProfileOn, true)
	}
}
