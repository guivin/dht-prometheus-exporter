package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	path            string
	name            string
	gpio            string
	maxRetries      int
	listenPort      int
	defaultLogLevel string
	temperatureUnit string
}

func ReadConfig() *Config {
	/**
	Reads YAML configuration file and create a struct containing the values
	**/
	viper.SetConfigName("dht-prometheus-exporter")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Sprintf("Error when reading config file: %v", err))
	}

	config := &Config{
		path:            viper.ConfigFileUsed(),
		name:            viper.GetString("name"),
		gpio:            fmt.Sprintf("GPIO%d", viper.GetInt("gpio_pin")),
		maxRetries:      viper.GetInt("max_retries"),
		listenPort:      viper.GetInt("listen_port"),
		defaultLogLevel: viper.GetString("log_level"),
		temperatureUnit: viper.GetString("temperature_unit"),
	}

	return config
}
