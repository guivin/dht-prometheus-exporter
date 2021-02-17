package main

import (
	"fmt"
	"github.com/MichaelS11/go-dht"
)

type Sensor struct {
	config            *Config
	temperatureSymbol string
	client            *dht.DHT
}

const CelsiusSymbol = "C"
const FahrenheitSymbol = "F"

func newSensor(config *Config) *Sensor {
	/**
	Sensor constructor
	*/
	var err error
	var client *dht.DHT
	var temperatureSymbol string

	lg.Info("Initializing the DHT22/AM2302 sensor on the host")
	err = dht.HostInit()
	if err != nil {
		lg.Panic("Failed to initialized DHT22/AM2302 sensor on the host: ", err)
	}

	if config.temperatureUnit == "celsius" {
		client, err = dht.NewDHT(config.gpio, dht.Celsius, "")
		temperatureSymbol = CelsiusSymbol
	} else {
		client, err = dht.NewDHT(config.gpio, dht.Fahrenheit, "")
		temperatureSymbol = FahrenheitSymbol
	}

	if err != nil {
		lg.Panic("Failed to create new DHT client: ", err)
	}

	return &Sensor{
		config:            config,
		temperatureSymbol: temperatureSymbol,
		client:            client,
	}
}

func (s *Sensor) readRetry() (humidity float64, temperature float64, err error) {
	/**
	Reads the sensor data with retry
	*/
	humidity, temperature, err = s.client.ReadRetry(s.config.maxRetries)
	if err != nil {
		lg.Error("Cannot retrieve humidity and temperature from the sensor: ", err)
	}
	lg.Info(fmt.Sprintf("Retrieved humidity=%f, temperature=%f°%s from the sensor",
		humidity, temperature, s.temperatureSymbol))
	return humidity, temperature, err
}
