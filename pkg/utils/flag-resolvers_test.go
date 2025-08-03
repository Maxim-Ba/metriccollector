package utils

import (
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
			name:      "env takes precedence",
			envValue:  "env_value",
			flag:      FlagValue[string]{Passed: true, Value: "flag_value"},
			fileValue: "file_value",
			expected:  "env_value",
		},
		{
			name:      "flag takes precedence when env is empty",
			envValue:  "",
			flag:      FlagValue[string]{Passed: true, Value: "flag_value"},
			fileValue: "file_value",
			expected:  "flag_value",
		},
		{
			name:      "file value used when env empty and flag not passed",
			envValue:  "",
			flag:      FlagValue[string]{Passed: false, Value: "flag_default"},
			fileValue: "file_value",
			expected:  "file_value",
		},
		{
			name:      "default flag value used when nothing else set",
			envValue:  "",
			flag:      FlagValue[string]{Passed: false, Value: "flag_default"},
			fileValue: "",
			expected:  "flag_default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolveString(tt.envValue, tt.flag, tt.fileValue)
			if result != tt.expected {
				t.Errorf("ResolveString() = %v, want %v", result, tt.expected)
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
			name:      "env takes precedence",
			envValue:  42,
			flag:      FlagValue[int]{Passed: true, Value: 10},
			fileValue: 20,
			expected:  42,
		},
		{
			name:      "flag takes precedence when env is zero",
			envValue:  0,
			flag:      FlagValue[int]{Passed: true, Value: 10},
			fileValue: 20,
			expected:  10,
		},
		{
			name:      "file value used when env zero and flag not passed",
			envValue:  0,
			flag:      FlagValue[int]{Passed: false, Value: 1},
			fileValue: 20,
			expected:  20,
		},
		{
			name:      "default flag value used when nothing else set",
			envValue:  0,
			flag:      FlagValue[int]{Passed: false, Value: 1},
			fileValue: 0,
			expected:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolveInt(tt.envValue, tt.flag, tt.fileValue)
			if result != tt.expected {
				t.Errorf("ResolveInt() = %v, want %v", result, tt.expected)
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
			name:      "env takes precedence when set",
			isEnvSet:  true,
			envValue:  true,
			flag:      FlagValue[bool]{Passed: true, Value: false},
			fileValue: false,
			expected:  true,
		},
		{
			name:      "flag takes precedence when env not set",
			isEnvSet:  false,
			envValue:  true, // should be ignored
			flag:      FlagValue[bool]{Passed: true, Value: true},
			fileValue: false,
			expected:  true,
		},
		{
			name:      "file value used when env not set and flag not passed",
			isEnvSet:  false,
			envValue:  false,
			flag:      FlagValue[bool]{Passed: false, Value: false},
			fileValue: true,
			expected:  true,
		},
		{
			name:      "default flag value used when nothing else set",
			isEnvSet:  false,
			envValue:  false,
			flag:      FlagValue[bool]{Passed: false, Value: true},
			fileValue: false,
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolveBool(tt.isEnvSet, tt.envValue, tt.flag, tt.fileValue)
			if result != tt.expected {
				t.Errorf("ResolveBool() = %v, want %v", result, tt.expected)
			}
		})
	}
}
