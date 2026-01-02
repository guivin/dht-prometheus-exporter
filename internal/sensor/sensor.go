package sensor

import (
	"fmt"

	"github.com/MichaelS11/go-dht"
	log "github.com/sirupsen/logrus"

	"github.com/guivin/dht-prometheus-exporter/internal/config"
)

const (
	// CelsiusSymbol is the symbol for Celsius temperature unit
	CelsiusSymbol = "C"
	// FahrenheitSymbol is the symbol for Fahrenheit temperature unit
	FahrenheitSymbol = "F"
)

// Reader defines the interface for reading sensor data.
// This interface allows for easy mocking in tests.
type Reader interface {
	// ReadData reads humidity and temperature from the sensor.
	// Returns an error if the sensor read fails.
	ReadData() (humidity, temperature float64, err error)

	// TemperatureUnit returns the temperature unit symbol ("C" or "F").
	TemperatureUnit() string

	// Name returns the sensor name.
	Name() string
}

// DHT22Sensor implements the Reader interface for DHT22/AM2302 sensors.
type DHT22Sensor struct {
	name              string
	maxRetries        int
	temperatureSymbol string
	client            *dht.DHT
	logger            *log.Logger
}

// HostInit initializes the DHT host. Must be called once before creating sensors.
func HostInit() error {
	return dht.HostInit()
}

// New creates a new DHT22 sensor reader.
// Returns an error if the sensor cannot be initialized.
func New(cfg *config.SensorConfig, logger *log.Logger) (*DHT22Sensor, error) {
	logger.Infof("Initializing DHT22/AM2302 sensor '%s' on %s", cfg.Name, cfg.GPIO)

	var client *dht.DHT
	var temperatureSymbol string
	var err error

	if cfg.TemperatureUnit == "celsius" {
		client, err = dht.NewDHT(cfg.GPIO, dht.Celsius, "")
		temperatureSymbol = CelsiusSymbol
	} else {
		client, err = dht.NewDHT(cfg.GPIO, dht.Fahrenheit, "")
		temperatureSymbol = FahrenheitSymbol
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create DHT client for sensor '%s': %w", cfg.Name, err)
	}

	return &DHT22Sensor{
		name:              cfg.Name,
		maxRetries:        cfg.MaxRetries,
		temperatureSymbol: temperatureSymbol,
		client:            client,
		logger:            logger,
	}, nil
}

// ReadData reads humidity and temperature from the sensor with retry logic.
// Returns an error if all retry attempts fail.
func (s *DHT22Sensor) ReadData() (humidity, temperature float64, err error) {
	humidity, temperature, err = s.client.ReadRetry(s.maxRetries)
	if err != nil {
		s.logger.Errorf("Cannot retrieve humidity and temperature from sensor '%s': %v", s.name, err)
		return 0, 0, err
	}

	s.logger.Infof("Retrieved humidity=%.2f%%, temperature=%.2f°%s from sensor '%s'",
		humidity, temperature, s.temperatureSymbol, s.name)

	return humidity, temperature, nil
}

// TemperatureUnit returns the temperature unit symbol for this sensor.
func (s *DHT22Sensor) TemperatureUnit() string {
	return s.temperatureSymbol
}

// Name returns the sensor name.
func (s *DHT22Sensor) Name() string {
	return s.name
}
