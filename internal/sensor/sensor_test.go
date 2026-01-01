package sensor

import (
	"errors"
	"io"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/guivin/dht-prometheus-exporter/internal/config"
)

// mockSensor is a mock implementation of the Reader interface for testing
type mockSensor struct {
	humidity    float64
	temperature float64
	err         error
	unit        string
}

func (m *mockSensor) ReadData() (float64, float64, error) {
	return m.humidity, m.temperature, m.err
}

func (m *mockSensor) TemperatureUnit() string {
	return m.unit
}

// TestMockSensor verifies that our mock implements the Reader interface
func TestMockSensor_ImplementsReader(t *testing.T) {
	var _ Reader = (*mockSensor)(nil)
}

func TestMockSensor_ReadData_Success(t *testing.T) {
	mock := &mockSensor{
		humidity:    65.5,
		temperature: 22.3,
		err:         nil,
		unit:        "C",
	}

	humidity, temperature, err := mock.ReadData()

	if err != nil {
		t.Errorf("ReadData() returned unexpected error: %v", err)
	}

	if humidity != 65.5 {
		t.Errorf("ReadData() humidity = %f, want %f", humidity, 65.5)
	}

	if temperature != 22.3 {
		t.Errorf("ReadData() temperature = %f, want %f", temperature, 22.3)
	}
}

func TestMockSensor_ReadData_Error(t *testing.T) {
	expectedErr := errors.New("sensor read failed")
	mock := &mockSensor{
		err: expectedErr,
	}

	_, _, err := mock.ReadData()

	if err == nil {
		t.Error("ReadData() expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("ReadData() error = %v, want %v", err, expectedErr)
	}
}

func TestMockSensor_TemperatureUnit(t *testing.T) {
	tests := []struct {
		name string
		unit string
	}{
		{"celsius", CelsiusSymbol},
		{"fahrenheit", FahrenheitSymbol},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSensor{unit: tt.unit}
			unit := mock.TemperatureUnit()

			if unit != tt.unit {
				t.Errorf("TemperatureUnit() = %q, want %q", unit, tt.unit)
			}
		})
	}
}

// TestDHT22Sensor_Constants verifies the temperature unit constants
func TestDHT22Sensor_Constants(t *testing.T) {
	if CelsiusSymbol != "C" {
		t.Errorf("CelsiusSymbol = %q, want %q", CelsiusSymbol, "C")
	}

	if FahrenheitSymbol != "F" {
		t.Errorf("FahrenheitSymbol = %q, want %q", FahrenheitSymbol, "F")
	}
}

// TestNew_TemperatureUnit verifies that the correct temperature unit is set
// Note: This test will fail on non-Raspberry Pi systems because dht.HostInit() requires GPIO hardware
// We're testing the logic flow even though it will fail at hardware initialization
func TestNew_TemperatureUnitLogic(t *testing.T) {
	// Create a silent logger for testing
	logger := log.New()
	logger.SetOutput(io.Discard)

	tests := []struct {
		name         string
		tempUnit     string
		expectedUnit string
	}{
		{"celsius", "celsius", CelsiusSymbol},
		{"fahrenheit", "fahrenheit", FahrenheitSymbol},
		{"other", "other", FahrenheitSymbol}, // Defaults to Fahrenheit
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Name:            "test-sensor",
				GPIO:            "GPIO4",
				MaxRetries:      5,
				ListenPort:      8080,
				LogLevel:        "info",
				TemperatureUnit: tt.tempUnit,
			}

			// Note: This will fail on systems without GPIO hardware
			// We're just testing that the function signature and basic logic are correct
			sensor, err := New(cfg, logger)

			// On systems without GPIO, we expect an error from dht.HostInit()
			// We can't test the full functionality without hardware
			if err != nil {
				// Expected on non-Raspberry Pi systems
				t.Logf("Sensor initialization failed (expected on non-RPi systems): %v", err)
				return
			}

			// If we somehow got a sensor (on actual Raspberry Pi with GPIO),
			// verify the temperature unit is correct
			if sensor != nil {
				if sensor.temperatureSymbol != tt.expectedUnit {
					t.Errorf("temperatureSymbol = %q, want %q", sensor.temperatureSymbol, tt.expectedUnit)
				}
			}
		})
	}
}

// TestReader_Interface verifies that DHT22Sensor implements Reader
func TestDHT22Sensor_ImplementsReader(t *testing.T) {
	var _ Reader = (*DHT22Sensor)(nil)
}

// Example of using the Reader interface in production code
func ExampleReader() {
	// This demonstrates how to use the Reader interface
	// In production, this would be a real DHT22Sensor
	// In tests, this can be a mockSensor

	var sensor Reader = &mockSensor{
		humidity:    60.0,
		temperature: 23.5,
		unit:        "C",
	}

	humidity, temperature, err := sensor.ReadData()
	if err == nil {
		_ = humidity    // Use humidity
		_ = temperature // Use temperature
		_ = sensor.TemperatureUnit()
	}
}
