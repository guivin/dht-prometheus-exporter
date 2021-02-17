package main

import "github.com/prometheus/client_golang/prometheus"

type Collector struct {
	sensor            *Sensor
	temperatureMetric *prometheus.Desc
	humidityMetric    *prometheus.Desc
}

func newCollector(s *Sensor) *Collector {
	lg.Debug("Creating a new prometheus collector for the sensor")
	return &Collector{
		sensor: s,
		temperatureMetric: prometheus.NewDesc("temperature",
			"Humidity measured by the sensor",
			nil, nil,
		),
		humidityMetric: prometheus.NewDesc("humidity",
			"Temperature measured by the sensor",
			nil, nil,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.temperatureMetric
	ch <- c.humidityMetric
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	humidity, temperature, err := c.sensor.readRetry()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.temperatureMetric, prometheus.CounterValue, temperature)
		ch <- prometheus.MustNewConstMetric(c.humidityMetric, prometheus.CounterValue, humidity)
	}
}
