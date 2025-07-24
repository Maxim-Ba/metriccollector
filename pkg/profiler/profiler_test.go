package profiler

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewProfiler_Disabled(t *testing.T) {
	p, err := New(false, "cpu.prof", "mem.prof")
	if err != nil {
		t.Fatalf("New with disabled profiler returned error: %v", err)
	}

	if p.isOn {
		t.Error("Expected profiler to be disabled")
	}
	if p.fileCPU != nil || p.fileMem != nil {
		t.Error("Expected file handles to be nil when profiler is disabled")
	}
}

func TestNewProfiler_Enabled(t *testing.T) {
	cpuProfile := "test_cpu.prof"
	memProfile := "test_mem.prof"

	// Cleanup before test
	err := os.RemoveAll("profiles")
	if err != nil {
		require.NoError(t, err)
	}

	p, err := New(true, cpuProfile, memProfile)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	defer func() {
		err = p.Close()
		require.NoError(t, err)
	}()

	if !p.isOn {
		t.Error("Expected profiler to be enabled")
	}
	if p.fileCPU == nil || p.fileMem == nil {
		t.Error("Expected file handles to be initialized")
	}

	// Verify files were created
	cpuPath := filepath.Join("profiles", cpuProfile)
	memPath := filepath.Join("profiles", memProfile)
	if _, err := os.Stat(cpuPath); os.IsNotExist(err) {
		t.Errorf("CPU profile file was not created: %v", err)
	}
	if _, err := os.Stat(memPath); os.IsNotExist(err) {
		t.Errorf("Memory profile file was not created: %v", err)
	}
}

func TestProfiler_StartStop(t *testing.T) {
	cpuProfile := "startstop_cpu.prof"
	memProfile := "startstop_mem.prof"

	// Cleanup before test
	err := os.RemoveAll("profiles")
	if err != nil {
		require.NoError(t, err)
	}
	p, err := New(true, cpuProfile, memProfile)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	defer func() {
		err = p.Close()
		require.NoError(t, err)
	}()
	// Test Start
	p.Start()

	// Give some time for the profiler to collect data
	time.Sleep(100 * time.Millisecond)

	// Test stop (via timer)
	// Wait for the timer to fire (30s is too long for tests, so we call stop directly)
	p.stop()

	// Verify files have content
	cpuPath := filepath.Join("profiles", cpuProfile)
	memPath := filepath.Join("profiles", memProfile)

	cpuInfo, err := os.Stat(cpuPath)
	if err != nil {
		t.Errorf("Could not stat CPU profile file: %v", err)
	}
	if cpuInfo.Size() == 0 {
		t.Error("CPU profile file is empty")
	}

	memInfo, err := os.Stat(memPath)
	if err != nil {
		t.Errorf("Could not stat memory profile file: %v", err)
	}
	if memInfo.Size() == 0 {
		t.Error("Memory profile file is empty")
	}
}

func TestProfiler_Close(t *testing.T) {
	cpuProfile := "close_cpu.prof"
	memProfile := "close_mem.prof"

	p, err := New(true, cpuProfile, memProfile)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if err = p.Close(); err != nil {
		t.Errorf("Close returned error: %v", err)
	}

	// Verify files are closed by trying to write to them
	_, err = p.fileCPU.WriteString("test")
	if err == nil {
		t.Error("Expected error when writing to closed CPU file")
	}

	_, err = p.fileMem.WriteString("test")
	if err == nil {
		t.Error("Expected error when writing to closed memory file")
	}
}

func TestProfiler_Close_Disabled(t *testing.T) {
	p := &Profiler{isOn: false}
	if err := p.Close(); err != nil {
		t.Errorf("Close on disabled profiler returned error: %v", err)
	}
}
