package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

// resetViper resets viper state between tests
func resetViper() {
	viper.Reset()
}

func TestLoad_ValidConfig(t *testing.T) {
	// Change to testdata directory for this test
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)
	defer resetViper()

	// Navigate to project root where testdata is
	testdataPath := filepath.Join(originalWd, "..", "..", "testdata")
	if err := os.Chdir(testdataPath); err != nil {
		t.Fatalf("Failed to change to testdata directory: %v", err)
	}

	// Rename test-config.yml to dht-prometheus-exporter.yml temporarily
	testConfigPath := "test-config.yml"
	expectedConfigPath := "dht-prometheus-exporter.yml"

	// If dht-prometheus-exporter.yml exists in testdata, back it up
	_, err = os.Stat(expectedConfigPath)
	backupNeeded := err == nil
	if backupNeeded {
		if err := os.Rename(expectedConfigPath, expectedConfigPath+".bak"); err != nil {
			t.Fatalf("Failed to backup existing config: %v", err)
		}
		defer os.Rename(expectedConfigPath+".bak", expectedConfigPath)
	}

	// Copy test-config.yml to dht-prometheus-exporter.yml
	if err := copyFile(testConfigPath, expectedConfigPath); err != nil {
		t.Fatalf("Failed to copy test config: %v", err)
	}
	defer os.Remove(expectedConfigPath)

	config, err := Load()
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}

	if config == nil {
		t.Fatal("Load() returned nil config")
	}

	// Verify config values match test-config.yml
	if config.Name != "test-sensor" {
		t.Errorf("Config.Name = %q, want %q", config.Name, "test-sensor")
	}

	if config.GPIO != "GPIO4" {
		t.Errorf("Config.GPIO = %q, want %q", config.GPIO, "GPIO4")
	}

	if config.MaxRetries != 5 {
		t.Errorf("Config.MaxRetries = %d, want %d", config.MaxRetries, 5)
	}

	if config.ListenPort != 9999 {
		t.Errorf("Config.ListenPort = %d, want %d", config.ListenPort, 9999)
	}

	if config.LogLevel != "info" {
		t.Errorf("Config.LogLevel = %q, want %q", config.LogLevel, "info")
	}

	if config.TemperatureUnit != "celsius" {
		t.Errorf("Config.TemperatureUnit = %q, want %q", config.TemperatureUnit, "celsius")
	}
}

func TestLoad_MissingConfig(t *testing.T) {
	defer resetViper()

	// Create a temporary empty directory
	tempDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory where no config exists
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	config, err := Load()
	if err == nil {
		t.Error("Load() expected error for missing config, got nil")
	}
	if config != nil {
		t.Errorf("Load() expected nil config on error, got %v", config)
	}
}

func TestLoad_GPIOPinFormatting(t *testing.T) {
	defer resetViper()

	// Create a temporary directory with a test config
	tempDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// Create a test config with a specific GPIO pin
	configContent := []byte(`---
name: test
gpio_pin: 17
max_retries: 10
listen_port: 8080
log_level: info
temperature_unit: celsius
`)
	configPath := filepath.Join(tempDir, "dht-prometheus-exporter.yml")
	if err := os.WriteFile(configPath, configContent, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	config, err := Load()
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}

	// Verify GPIO pin is formatted with "GPIO" prefix
	expectedGPIO := "GPIO17"
	if config.GPIO != expectedGPIO {
		t.Errorf("Config.GPIO = %q, want %q", config.GPIO, expectedGPIO)
	}
}

// copyFile is a helper function to copy a file
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
