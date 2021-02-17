package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"os"
)

type Collector struct {
	sensor            *Sensor
	temperatureMetric *prometheus.Desc
	humidityMetric    *prometheus.Desc
}

func newCollector(s *Sensor) *Collector {
	lg.Debug("Creating a new prometheus collector for the sensor")
	return &Collector{
		sensor: s,
		temperatureMetric: prometheus.NewDesc("dht_temperature_degree",
			"Temperature degree measured by the sensor",
			[]string{"dht_name", "hostname", "unit"}, nil,
		),
		humidityMetric: prometheus.NewDesc("dht_humidity_percent",
			"Humidity percent measured by the sensor",
			[]string{"dht_name", "hostname"}, nil,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.temperatureMetric
	ch <- c.humidityMetric
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	humidity, temperature, _ := c.sensor.readRetry()
	hostname, err := os.Hostname()
	temperatureUnit := c.sensor.config.temperatureUnit
	dhtName := c.sensor.config.name
	if err != nil {
		lg.Error("Failed to get hostname")
	}
	ch <- prometheus.MustNewConstMetric(c.temperatureMetric, prometheus.CounterValue, temperature, dhtName, hostname, temperatureUnit)
	ch <- prometheus.MustNewConstMetric(c.humidityMetric, prometheus.CounterValue, humidity, dhtName, hostname)
}
