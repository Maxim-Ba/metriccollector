package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseEnv(t *testing.T) {
	// Сохраняем текущее окружение и очищаем его для теста
	oldEnv := os.Environ()
	os.Clearenv()

	// Восстанавливаем окружение после теста
	defer func() {
		os.Clearenv()
		for _, env := range oldEnv {
			keyVal := strings.SplitN(env, "=", 2)
			if len(keyVal) != 2 {
				continue
			}
			if err := os.Setenv(keyVal[0], keyVal[1]); err != nil {
				t.Errorf("Failed to restore env %s: %v", keyVal[0], err)
			}
		}
	}()

	// Устанавливаем тестовые переменные окружения
	testVars := map[string]string{
		"ADDRESS":           ":8080",
		"STORE_INTERVAL":    "10",
		"FILE_STORAGE_PATH": "/tmp/test",
		"RESTORE":           "true",
		"LOG_LEVEL":         "debug",
		"DATABASE_DSN":      "postgres://user:pass@localhost:5432/db",
		"MIGRATIONS_PATH":   "./migrations",
		"KEY":               "secret",
		"CPU_FILE":          "cpu.prof",
		"MEM_FILE":          "mem.prof",
		"IS_PROFILE_ON":     "true",
	}

	for k, v := range testVars {
		err := os.Setenv(k, v)
		if err != nil {
			require.NoError(t, err)
		}
	}

	// Вызываем тестируемую функцию
	cfg := ParseEnv()

	// Проверяем результаты
	if cfg.Address != ":8080" {
		t.Errorf("expected Addres ':8080', got '%s'", cfg.Address)
	}
	if cfg.StoreIntervalSecond != 10 {
		t.Errorf("expected StoreIntervalSecond 10, got %d", cfg.StoreIntervalSecond)
	}
	if cfg.StoragePath != "/tmp/test" {
		t.Errorf("expected StoragePath '/tmp/test', got '%s'", cfg.StoragePath)
	}
	if !cfg.Restore {
		t.Error("expected Restore true, got false")
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected LogLevel 'debug', got '%s'", cfg.LogLevel)
	}
	if cfg.DatabaseDSN != "postgres://user:pass@localhost:5432/db" {
		t.Errorf("expected DatabaseDSN 'postgres://user:pass@localhost:5432/db', got '%s'", cfg.DatabaseDSN)
	}
	if cfg.MigrationsPath != "./migrations" {
		t.Errorf("expected MigrationsPath './migrations', got '%s'", cfg.MigrationsPath)
	}
	if cfg.Key != "secret" {
		t.Errorf("expected Key 'secret', got '%s'", cfg.Key)
	}
	if cfg.ProfileFileCPU != "cpu.prof" {
		t.Errorf("expected ProfileFileCPU 'cpu.prof', got '%s'", cfg.ProfileFileCPU)
	}
	if cfg.ProfileFileMem != "mem.prof" {
		t.Errorf("expected ProfileFileMem 'mem.prof', got '%s'", cfg.ProfileFileMem)
	}
	if !cfg.IsProfileOn {
		t.Error("expected IsProfileOn true, got false")
	}
}

func TestIsRestoreSet(t *testing.T) {
	// Сохраняем текущее окружение и очищаем его для теста
	oldEnv := os.Environ()
	os.Clearenv()

	// Восстанавливаем окружение после теста
	defer func() {
		os.Clearenv()
		for _, env := range oldEnv {
			keyVal := strings.SplitN(env, "=", 2)
			if len(keyVal) != 2 {
				continue
			}
			if err := os.Setenv(keyVal[0], keyVal[1]); err != nil {
				t.Errorf("Failed to restore env %s: %v", keyVal[0], err)
			}
		}
	}()

	tests := []struct {
		name     string
		setValue string
		want     bool
	}{
		{"not set", "", false},
		{"set empty", "", true},
		{"set with value", "true", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setValue != "" || tt.name == "set empty" {
				err := os.Setenv("RESTORE", tt.setValue)
				if err != nil {
					t.Fatalf("Failed to restore env: %v", err)
				}
			}

			if got := isRestoreSet(); got != tt.want {
				t.Errorf("isRestoreSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIntervalSet(t *testing.T) {
	// Сохраняем текущее окружение и очищаем его для теста
	oldEnv := os.Environ()
	os.Clearenv()

	// Восстанавливаем окружение после теста
	defer func() {
		os.Clearenv()
		for _, env := range oldEnv {
			keyVal := strings.SplitN(env, "=", 2)
			if len(keyVal) != 2 {
				continue
			}
			if err := os.Setenv(keyVal[0], keyVal[1]); err != nil {
				t.Errorf("Failed to restore env %s: %v", keyVal[0], err)
			}
		}
	}()

	tests := []struct {
		name     string
		setValue string
		want     bool
	}{
		{"not set", "", false},
		{"set empty", "", true},
		{"set with value", "10", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setValue != "" || tt.name == "set empty" {
				err := os.Setenv("STORE_INTERVAL", tt.setValue)
				if err != nil {
					t.Fatalf("Failed to restore env: %v", err)
				}
			}

			if got := isIntervalSet(); got != tt.want {
				t.Errorf("isIntervalSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMigrationsPathSet(t *testing.T) {
	// Сохраняем текущее окружение и очищаем его для теста
	oldEnv := os.Environ()
	os.Clearenv()

	// Восстанавливаем окружение после теста
	defer func() {
		os.Clearenv()
		for _, env := range oldEnv {
			keyVal := strings.SplitN(env, "=", 2)
			if len(keyVal) != 2 {
				continue
			}
			if err := os.Setenv(keyVal[0], keyVal[1]); err != nil {
				t.Errorf("Failed to restore env %s: %v", keyVal[0], err)
			}
		}
	}()

	tests := []struct {
		name     string
		setValue string
		want     bool
	}{
		{"not set", "", false},
		{"set empty", "", true},
		{"set with value", "/path", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setValue != "" || tt.name == "set empty" {
				err := os.Setenv("MIGRATIONS_PATH", tt.setValue)
				if err != nil {
					t.Fatalf("Failed to restore env: %v", err)
				}
			}

			if got := isMigrationsPathSet(); got != tt.want {
				t.Errorf("isMigrationsPathSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsProfileOnSet(t *testing.T) {
	// Сохраняем текущее окружение и очищаем его для теста
	oldEnv := os.Environ()
	os.Clearenv()

	// Восстанавливаем окружение после теста
	defer func() {
		os.Clearenv()

		for _, env := range oldEnv {
			keyVal := strings.SplitN(env, "=", 2)
			if len(keyVal) != 2 {
				continue
			}
			if err := os.Setenv(keyVal[0], keyVal[1]); err != nil {
				t.Errorf("Failed to restore env %s: %v", keyVal[0], err)
			}
		}
	}()

	tests := []struct {
		name     string
		setValue string
		want     bool
	}{
		{"not set", "", false},
		{"set empty", "", true},
		{"set with value", "true", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setValue != "" || tt.name == "set empty" {
				err := os.Setenv("IS_PROFILE_ON", tt.setValue)
				if err != nil {
					t.Fatalf("Failed to restore env: %v", err)
				}
			}

			if got := isProfileOnSet(); got != tt.want {
				t.Errorf("isProfileOnSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
