package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// SensorConfig holds the configuration for a single DHT sensor.
type SensorConfig struct {
	Name            string
	GPIO            string
	MaxRetries      int
	TemperatureUnit string
}

// Config holds the application configuration loaded from YAML file.
type Config struct {
	Sensors    []SensorConfig
	ListenPort int
	LogLevel   string
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

	var sensors []SensorConfig
	sensorsRaw := viper.Get("sensors")
	if sensorsRaw != nil {
		sensorsList, ok := sensorsRaw.([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid sensors configuration format")
		}
		for i, s := range sensorsList {
			sensorMap, ok := s.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid sensor configuration at index %d", i)
			}
			sensor := SensorConfig{
				Name:            getString(sensorMap, "name"),
				GPIO:            fmt.Sprintf("GPIO%d", getInt(sensorMap, "gpio_pin")),
				MaxRetries:      getInt(sensorMap, "max_retries"),
				TemperatureUnit: getString(sensorMap, "temperature_unit"),
			}
			sensors = append(sensors, sensor)
		}
	}

	if len(sensors) == 0 {
		return nil, fmt.Errorf("no sensors configured")
	}

	config := &Config{
		Sensors:    sensors,
		ListenPort: viper.GetInt("listen_port"),
		LogLevel:   viper.GetString("log_level"),
	}

	return config, nil
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case int:
			return n
		case float64:
			return int(n)
		}
	}
	return 0
}
