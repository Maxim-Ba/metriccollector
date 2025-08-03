package config

import (
	"os"
	"testing"

	"github.com/Maxim-Ba/metriccollector/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestResolveString(t *testing.T) {
	tests := []struct {
		name      string
		envValue  string
		flag      utils.FlagValue[string]
		fileValue string
		expected  string
	}{
		{
			name:      "env value takes precedence",
			envValue:  "env_value",
			flag:      utils.FlagValue[string]{Value: "flag_value", Passed: true},
			fileValue: "file_value",
			expected:  "env_value",
		},
		{
			name:      "flag value when env empty and flag passed",
			envValue:  "",
			flag:      utils.FlagValue[string]{Value: "flag_value", Passed: true},
			fileValue: "file_value",
			expected:  "flag_value",
		},
		{
			name:      "file value when env empty and flag not passed",
			envValue:  "",
			flag:      utils.FlagValue[string]{Value: "flag_value", Passed: false},
			fileValue: "file_value",
			expected:  "file_value",
		},
		{
			name:      "default flag value when all empty",
			envValue:  "",
			flag:      utils.FlagValue[string]{Value: "default_flag_value", Passed: false},
			fileValue: "",
			expected:  "default_flag_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ResolveString(tt.envValue, tt.flag, tt.fileValue)
			if result != tt.expected {
				t.Errorf("resolveString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestResolveInt(t *testing.T) {
	tests := []struct {
		name      string
		envValue  int
		flag      utils.FlagValue[int]
		fileValue int
		expected  int
	}{
		{
			name:      "env value takes precedence",
			envValue:  42,
			flag:      utils.FlagValue[int]{Value: 10, Passed: true},
			fileValue: 20,
			expected:  42,
		},
		{
			name:      "flag value when env empty and flag passed",
			envValue:  0,
			flag:      utils.FlagValue[int]{Value: 10, Passed: true},
			fileValue: 20,
			expected:  10,
		},
		{
			name:      "file value when env empty and flag not passed",
			envValue:  0,
			flag:      utils.FlagValue[int]{Value: 10, Passed: false},
			fileValue: 20,
			expected:  20,
		},
		{
			name:      "default flag value when all empty",
			envValue:  0,
			flag:      utils.FlagValue[int]{Value: 5, Passed: false},
			fileValue: 0,
			expected:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ResolveInt(tt.envValue, tt.flag, tt.fileValue)
			if result != tt.expected {
				t.Errorf("resolveInt() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestResolveBool(t *testing.T) {
	tests := []struct {
		name      string
		isEnvSet  bool
		envValue  bool
		flag      utils.FlagValue[bool]
		fileValue bool
		expected  bool
	}{
		{
			name:      "env value takes precedence when set",
			isEnvSet:  true,
			envValue:  true,
			flag:      utils.FlagValue[bool]{Value: false, Passed: true},
			fileValue: false,
			expected:  true,
		},
		{
			name:      "flag value when env not set and flag passed",
			isEnvSet:  false,
			envValue:  false,
			flag:      utils.FlagValue[bool]{Value: true, Passed: true},
			fileValue: false,
			expected:  true,
		},
		{
			name:      "file value when env not set and flag not passed",
			isEnvSet:  false,
			envValue:  false,
			flag:      utils.FlagValue[bool]{Value: false, Passed: false},
			fileValue: true,
			expected:  true,
		},
		{
			name:      "default flag value when all empty",
			isEnvSet:  false,
			envValue:  false,
			flag:      utils.FlagValue[bool]{Value: true, Passed: false},
			fileValue: false,
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ResolveBool(tt.isEnvSet, tt.envValue, tt.flag, tt.fileValue)
			if result != tt.expected {
				t.Errorf("resolveBool() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetParamsByConfigPath(t *testing.T) {
	// Создаем временный файл конфигурации
	configContent := `{
		"address": "127.0.0.1:8080",
		"store_interval": 10,
		"store_file": "/tmp/test.json",
		"restore": true,
		"log_level": "debug"
	}`

	tmpFile, err := os.CreateTemp("", "config_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	tests := []struct {
		name       string
		configPath string
		wantErr    bool
		expected   Parameters
	}{
		{
			name:       "valid config file",
			configPath: tmpFile.Name(),
			wantErr:    false,
			expected: Parameters{
				Address:             "127.0.0.1:8080",
				StoreIntervalSecond: 10,
				StoragePath:         "/tmp/test.json",
				Restore:             true,
				LogLevel:            "debug",
			},
		},
		{
			name:       "empty path",
			configPath: "",
			wantErr:    false,
			expected:   Parameters{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getParamsByConfigPath(tt.configPath)
			assert.NoError(t, err)

			if tt.wantErr {
				// Проверяем, что вернулся пустой Parameters при ошибке
				if got != tt.expected {
					t.Errorf("getParamsByConfigPath() = %v, want %v", got, tt.expected)
				}
			} else {
				// Проверяем конкретные поля
				if got.Address != tt.expected.Address ||
					got.StoreIntervalSecond != tt.expected.StoreIntervalSecond ||
					got.StoragePath != tt.expected.StoragePath ||
					got.Restore != tt.expected.Restore ||
					got.LogLevel != tt.expected.LogLevel {
					t.Errorf("getParamsByConfigPath() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	// Этот тест сложнее, так как зависит от флагов и переменных окружения
	// Можно использовать моки или тестовые переменные окружения

	// Сохраняем текущие переменные окружения
	oldEnv := os.Environ()
	defer func() {
		// Восстанавливаем переменные окружения после теста
		os.Clearenv()
		for _, env := range oldEnv {
			os.Setenv(env, "")
		}
	}()

	// Устанавливаем тестовые переменные окружения
	os.Setenv("ADDRESS", "env_address:8080")
	os.Setenv("STORE_INTERVAL", "30")
	os.Setenv("FILE_STORAGE_PATH", "/env/path.json")
	os.Setenv("RESTORE", "true")
	os.Setenv("LOG_LEVEL", "info")

	// Создаем временный файл конфигурации
	configContent := `{
		"address": "file_address:9090",
		"store_interval": 20,
		"store_file": "/file/path.json",
		"restore": false,
		"log_level": "debug"
	}`

	tmpFile, err := os.CreateTemp("", "config_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Сохраняем оригинальные флаги и восстанавливаем после теста
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Устанавливаем тестовые флаги
	os.Args = []string{
		"cmd",
		"-a", "flag_address:7070",
		"-i", "10",
		"-f", "/flag/path.json",
		"-r",
		"-l", "warn",
		"-c", tmpFile.Name(),
	}

	params := New()

	// Проверяем, что переменные окружения имеют приоритет
	if params.Address != "env_address:8080" {
		t.Errorf("New() Address = %v, want env_address:8080", params.Address)
	}
	if params.StoreIntervalSecond != 30 {
		t.Errorf("New() StoreIntervalSecond = %v, want 30", params.StoreIntervalSecond)
	}
	if params.StoragePath != "/env/path.json" {
		t.Errorf("New() StoragePath = %v, want /env/path.json", params.StoragePath)
	}
	if !params.Restore {
		t.Errorf("New() Restore = %v, want true", params.Restore)
	}
	if params.LogLevel != "info" {
		t.Errorf("New() LogLevel = %v, want info", params.LogLevel)
	}
}
