package collector

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/guivin/dht-prometheus-exporter/internal/sensor"
)

// Collector implements the prometheus.Collector interface for DHT sensor metrics.
type Collector struct {
	sensor            sensor.Reader
	logger            *log.Logger
	hostname          string
	temperatureMetric *prometheus.Desc
	humidityMetric    *prometheus.Desc
}

// New creates a new Collector for the given sensor.
// The hostname is retrieved once during initialization to avoid repeated lookups.
func New(s sensor.Reader, logger *log.Logger) *Collector {
	logger.WithField("sensor", s.Name()).Debug("Creating Prometheus collector")

	hostname, err := os.Hostname()
	if err != nil {
		logger.WithError(err).Warn("Failed to get hostname, using empty string")
		hostname = ""
	}

	return &Collector{
		sensor:   s,
		logger:   logger,
		hostname: hostname,
		temperatureMetric: prometheus.NewDesc(
			"dht_temperature_degree",
			"Temperature degree measured by the sensor",
			[]string{"dht_name", "hostname", "gpio", "unit"}, nil,
		),
		humidityMetric: prometheus.NewDesc(
			"dht_humidity_percent",
			"Humidity percent measured by the sensor",
			[]string{"dht_name", "hostname", "gpio"}, nil,
		),
	}
}

// Describe sends the descriptors of the metrics to the provided channel.
// This is required by the prometheus.Collector interface.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.temperatureMetric
	ch <- c.humidityMetric
}

// Collect reads sensor data and sends metrics to the provided channel.
// If sensor reading fails, no metrics are emitted.
// CRITICAL FIX: Uses GaugeValue instead of CounterValue (temperature/humidity are gauges, not counters)
// CRITICAL FIX: Checks and handles sensor read errors instead of ignoring them
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	humidity, temperature, err := c.sensor.ReadData()
	if err != nil {
		// Error already logged by sensor.ReadData(), just skip metric collection
		return
	}

	temperatureUnit := c.sensor.TemperatureUnit()

	// CRITICAL FIX: Use GaugeValue instead of CounterValue
	// Temperature and humidity are gauge metrics (can go up or down), not counters
	ch <- prometheus.MustNewConstMetric(
		c.temperatureMetric,
		prometheus.GaugeValue, // Changed from CounterValue
		temperature,
		c.sensor.Name(),
		c.hostname,
		c.sensor.GPIO(),
		temperatureUnit,
	)

	ch <- prometheus.MustNewConstMetric(
		c.humidityMetric,
		prometheus.GaugeValue, // Changed from CounterValue
		humidity,
		c.sensor.Name(),
		c.hostname,
		c.sensor.GPIO(),
	)
}
