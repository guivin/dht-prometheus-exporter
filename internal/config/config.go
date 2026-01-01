package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds the application configuration loaded from YAML file.
type Config struct {
	Name            string
	GPIO            string
	MaxRetries      int
	ListenPort      int
	LogLevel        string
	TemperatureUnit string
}

// Load reads and validates the configuration from the default locations.
// It searches for the config file in /etc, $HOME, and current directory.
// Returns an error if the config file cannot be read or parsed.
func Load() (*Config, error) {
	viper.SetConfigName("dht-prometheus-exporter")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &Config{
		Name:            viper.GetString("name"),
		GPIO:            fmt.Sprintf("GPIO%d", viper.GetInt("gpio_pin")),
		MaxRetries:      viper.GetInt("max_retries"),
		ListenPort:      viper.GetInt("listen_port"),
		LogLevel:        viper.GetString("log_level"),
		TemperatureUnit: viper.GetString("temperature_unit"),
	}

	return config, nil
}
