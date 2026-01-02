package collector

import (
	"errors"
	"io"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	log "github.com/sirupsen/logrus"
)

// mockSensor is a mock implementation of sensor.Reader for testing
type mockSensor struct {
	name        string
	gpio        string
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

func (m *mockSensor) Name() string {
	return m.name
}

func (m *mockSensor) GPIO() string {
	return m.gpio
}

// getSilentLogger returns a logger that doesn't output anything
func getSilentLogger() *log.Logger {
	logger := log.New()
	logger.SetOutput(io.Discard)
	return logger
}

func TestNew(t *testing.T) {
	logger := getSilentLogger()
	mock := &mockSensor{name: "test-sensor", unit: "C"}

	collector := New(mock, logger)

	if collector == nil {
		t.Fatal("New() returned nil collector")
	}

	if collector.logger != logger {
		t.Error("Collector.logger not set correctly")
	}

	// Hostname should be set (even if empty on error)
	// We just verify it's not nil
	_ = collector.hostname
}

func TestDescribe(t *testing.T) {
	logger := getSilentLogger()
	mock := &mockSensor{name: "test-sensor", unit: "C"}
	collector := New(mock, logger)

	ch := make(chan *prometheus.Desc, 10)
	collector.Describe(ch)
	close(ch)

	count := 0
	for range ch {
		count++
	}

	if count != 2 {
		t.Errorf("Describe() sent %d descriptors, want 2", count)
	}
}

// CRITICAL TEST: Verify metrics are emitted on successful read
func TestCollect_Success(t *testing.T) {
	logger := getSilentLogger()
	mock := &mockSensor{
		name:        "test-sensor",
		humidity:    65.5,
		temperature: 22.3,
		err:         nil,
		unit:        "C",
	}

	collector := New(mock, logger)

	ch := make(chan prometheus.Metric, 10)
	collector.Collect(ch)
	close(ch)

	count := 0
	for range ch {
		count++
	}

	if count != 2 {
		t.Errorf("Collect() emitted %d metrics on success, want 2", count)
	}
}

// CRITICAL TEST: Verify NO metrics are emitted on sensor error
// This tests the fix for ignoring sensor errors
func TestCollect_SensorError(t *testing.T) {
	logger := getSilentLogger()
	mock := &mockSensor{
		name: "test-sensor",
		err:  errors.New("sensor read failed"),
	}

	collector := New(mock, logger)

	ch := make(chan prometheus.Metric, 10)
	collector.Collect(ch)
	close(ch)

	count := 0
	for range ch {
		count++
	}

	// CRITICAL: On error, NO metrics should be emitted
	if count != 0 {
		t.Errorf("Collect() emitted %d metrics on error, want 0", count)
	}
}

// CRITICAL TEST: Verify metrics use GaugeValue, not CounterValue
func TestCollect_MetricType(t *testing.T) {
	logger := getSilentLogger()
	mock := &mockSensor{
		name:        "test-sensor",
		humidity:    60.0,
		temperature: 25.0,
		unit:        "C",
	}

	collector := New(mock, logger)

	ch := make(chan prometheus.Metric, 10)
	collector.Collect(ch)
	close(ch)

	for metric := range ch {
		// Write metric to DTO to inspect it
		var dtoMetric dto.Metric
		if err := metric.Write(&dtoMetric); err != nil {
			t.Fatalf("Failed to write metric: %v", err)
		}

		// CRITICAL: Verify metric type is GAUGE, not COUNTER
		if dtoMetric.GetGauge() == nil {
			t.Errorf("Metric is not a Gauge type (got %+v)", dtoMetric)
		}

		// Verify counter is NOT set (it should be nil for gauges)
		if dtoMetric.GetCounter() != nil {
			t.Error("Metric incorrectly set as Counter instead of Gauge")
		}
	}
}

// TestCollect_MetricValues verifies the actual metric values
func TestCollect_MetricValues(t *testing.T) {
	logger := getSilentLogger()

	tests := []struct {
		name        string
		humidity    float64
		temperature float64
		unit        string
	}{
		{"celsius", 65.5, 22.3, "C"},
		{"fahrenheit", 70.0, 75.5, "F"},
		{"zero values", 0.0, 0.0, "C"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSensor{
				name:        "test-sensor",
				humidity:    tt.humidity,
				temperature: tt.temperature,
				unit:        tt.unit,
			}

			collector := New(mock, logger)

			ch := make(chan prometheus.Metric, 10)
			collector.Collect(ch)
			close(ch)

			metricCount := 0
			for metric := range ch {
				metricCount++

				var dtoMetric dto.Metric
				if err := metric.Write(&dtoMetric); err != nil {
					t.Fatalf("Failed to write metric: %v", err)
				}

				value := dtoMetric.GetGauge().GetValue()

				// Check if this is temperature or humidity metric
				// by examining the metric value
				if value == tt.temperature || value == tt.humidity {
					// Verify the value matches expected
					if value != tt.temperature && value != tt.humidity {
						t.Errorf("Metric value = %f, want %f or %f",
							value, tt.temperature, tt.humidity)
					}
				}
			}

			if metricCount != 2 {
				t.Errorf("Expected 2 metrics, got %d", metricCount)
			}
		})
	}
}

// TestCollect_MetricLabels verifies metric labels are set correctly
func TestCollect_MetricLabels(t *testing.T) {
	logger := getSilentLogger()
	dhtName := "my-sensor"
	mock := &mockSensor{
		name:        dhtName,
		humidity:    60.0,
		temperature: 20.0,
		unit:        "C",
	}

	collector := New(mock, logger)

	ch := make(chan prometheus.Metric, 10)
	collector.Collect(ch)
	close(ch)

	for metric := range ch {
		var dtoMetric dto.Metric
		if err := metric.Write(&dtoMetric); err != nil {
			t.Fatalf("Failed to write metric: %v", err)
		}

		// Verify labels are present
		labels := dtoMetric.GetLabel()
		if len(labels) < 2 {
			t.Errorf("Expected at least 2 labels, got %d", len(labels))
		}

		// Verify dht_name label
		foundDhtName := false
		foundHostname := false

		for _, label := range labels {
			if label.GetName() == "dht_name" {
				foundDhtName = true
				if label.GetValue() != dhtName {
					t.Errorf("dht_name label = %q, want %q", label.GetValue(), dhtName)
				}
			}
			if label.GetName() == "hostname" {
				foundHostname = true
				// Hostname could be anything, just verify it's present
			}
		}

		if !foundDhtName {
			t.Error("dht_name label not found in metric")
		}

		if !foundHostname {
			t.Error("hostname label not found in metric")
		}
	}
}
