package config

import (
	"os"
	"testing"
)



func TestResolveString(t *testing.T) {
	tests := []struct {
		name      string
		envValue  string
		flag      FlagValue[string]
		fileValue string
		expected  string
	}{
		{
			name:      "env value takes precedence",
			envValue:  "env_value",
			flag:      FlagValue[string]{Value: "flag_value", Passed: true},
			fileValue: "file_value",
			expected:  "env_value",
		},
		{
			name:      "flag value when env empty and flag passed",
			envValue:  "",
			flag:      FlagValue[string]{Value: "flag_value", Passed: true},
			fileValue: "file_value",
			expected:  "flag_value",
		},
		{
			name:      "file value when env empty and flag not passed",
			envValue:  "",
			flag:      FlagValue[string]{Value: "flag_value", Passed: false},
			fileValue: "file_value",
			expected:  "file_value",
		},
		{
			name:      "default flag value when all empty",
			envValue:  "",
			flag:      FlagValue[string]{Value: "default_flag_value", Passed: false},
			fileValue: "",
			expected:  "default_flag_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveString(tt.envValue, tt.flag, tt.fileValue)
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
		flag      FlagValue[int]
		fileValue int
		expected  int
	}{
		{
			name:      "env value takes precedence",
			envValue:  42,
			flag:      FlagValue[int]{Value: 10, Passed: true},
			fileValue: 20,
			expected:  42,
		},
		{
			name:      "flag value when env empty and flag passed",
			envValue:  0,
			flag:      FlagValue[int]{Value: 10, Passed: true},
			fileValue: 20,
			expected:  10,
		},
		{
			name:      "file value when env empty and flag not passed",
			envValue:  0,
			flag:      FlagValue[int]{Value: 10, Passed: false},
			fileValue: 20,
			expected:  20,
		},
		{
			name:      "default flag value when all empty",
			envValue:  0,
			flag:      FlagValue[int]{Value: 5, Passed: false},
			fileValue: 0,
			expected:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveInt(tt.envValue, tt.flag, tt.fileValue)
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
		flag      FlagValue[bool]
		fileValue bool
		expected  bool
	}{
		{
			name:      "env value takes precedence when set",
			isEnvSet:  true,
			envValue:  true,
			flag:      FlagValue[bool]{Value: false, Passed: true},
			fileValue: false,
			expected:  true,
		},
		{
			name:      "flag value when env not set and flag passed",
			isEnvSet:  false,
			envValue:  false,
			flag:      FlagValue[bool]{Value: true, Passed: true},
			fileValue: false,
			expected:  true,
		},
		{
			name:      "file value when env not set and flag not passed",
			isEnvSet:  false,
			envValue:  false,
			flag:      FlagValue[bool]{Value: false, Passed: false},
			fileValue: true,
			expected:  true,
		},
		{
			name:      "default flag value when all empty",
			isEnvSet:  false,
			envValue:  false,
			flag:      FlagValue[bool]{Value: true, Passed: false},
			fileValue: false,
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveBool(tt.isEnvSet, tt.envValue, tt.flag, tt.fileValue)
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
		"report_interval": 10,
		"poll_interval": 5,
		"log_level": "debug",
		"rate_limit": 100,
		"crypto_key": "/path/to/key"
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
				Address:        "127.0.0.1:8080",
				ReportInterval: 10,
				PollInterval:   5,
				LogLevel:       "debug",
				RateLimit:      100,
				CryptoKeyPath:  "/path/to/key",
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
			got := getParamsByConfigPath(tt.configPath)
			
			if tt.wantErr {
				// Проверяем, что вернулся пустой Parameters при ошибке
				if got != tt.expected {
					t.Errorf("getParamsByConfigPath() = %v, want %v", got, tt.expected)
				}
			} else {
				// Проверяем конкретные поля
				if got.Address != tt.expected.Address ||
					got.ReportInterval != tt.expected.ReportInterval ||
					got.PollInterval != tt.expected.PollInterval ||
					got.LogLevel != tt.expected.LogLevel ||
					got.RateLimit != tt.expected.RateLimit ||
					got.CryptoKeyPath != tt.expected.CryptoKeyPath {
					t.Errorf("getParamsByConfigPath() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
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
	os.Setenv("REPORT_INTERVAL", "30")
	os.Setenv("POLL_INTERVAL", "15")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("RATE_LIMIT", "200")
	os.Setenv("CRYPTO_KEY", "/env/key/path")

	// Создаем временный файл конфигурации
	configContent := `{
		"address": "file_address:9090",
		"report_interval": 20,
		"poll_interval": 10,
		"log_level": "debug",

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

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"cmd",
		"-a", "flag_address:7070",
		"-r", "5",
		"-p", "3",
		"-c", tmpFile.Name(),
	}

	params := New()

	if params.Address != "env_address:8080" {
		t.Errorf("New() Address = %v, want env_address:8080", params.Address)
	}
	if params.ReportInterval != 30 {
		t.Errorf("New() ReportInterval = %v, want 30", params.ReportInterval)
	}
	if params.PollInterval != 15 {
		t.Errorf("New() PollInterval = %v, want 15", params.PollInterval)
	}
	if params.LogLevel != "info" {
		t.Errorf("New() LogLevel = %v, want info", params.LogLevel)
	}
	if params.RateLimit != 200 {
		t.Errorf("New() RateLimit = %v, want 200", params.RateLimit)
	}
	if params.CryptoKeyPath != "/env/key/path" {
		t.Errorf("New() CryptoKeyPath = %v, want /env/key/path", params.CryptoKeyPath)
	}
}
