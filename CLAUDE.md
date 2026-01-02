# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Prometheus exporter for DHT22/AM2302 temperature and humidity sensors, designed to run on Raspberry Pi. It exposes sensor readings via HTTP endpoint for Prometheus scraping. Supports multiple sensors from a single Raspberry Pi.

## Build and Development Commands

### Build
```bash
make build           # Build the binary
make all             # Run tests + build
```

### Test
```bash
make test            # Run tests with race detection and coverage
make test-coverage   # Generate HTML coverage report
```

### Install/Uninstall
```bash
make install        # Copy binary to /usr/bin (requires sudo)
make uninstall      # Remove binary from /usr/bin (requires sudo)
```

### Cleanup
```bash
make clean          # Remove build artifacts and coverage files
make mod-tidy       # Tidy Go modules
```

## Architecture

### Core Components

- **cmd/dht-prometheus-exporter/main.go**: Entry point that orchestrates initialization and HTTP server setup
  - Reads config from YAML (searches /etc, $HOME, and current directory)
  - Initializes logger, then iterates over configured sensors
  - Creates a sensor and collector for each configured sensor
  - Starts HTTP server on configured port with /metrics endpoint

- **internal/config/config.go**: Configuration management using Viper
  - Supports YAML config file named `dht-prometheus-exporter.yml`
  - `SensorConfig` struct: name, gpio, max_retries, temperature_unit (per sensor)
  - `Config` struct: sensors array, listen_port, log_level (global)

- **internal/sensor/sensor.go**: DHT22/AM2302 sensor interface wrapper
  - Uses go-dht library for GPIO communication
  - Supports both Celsius and Fahrenheit
  - Provides retry mechanism for reliable sensor reads
  - `HostInit()` must be called once before creating sensor instances
  - `Reader` interface: ReadData(), TemperatureUnit(), Name()

- **internal/collector/collector.go**: Prometheus collector implementation
  - Implements prometheus.Collector interface (Describe/Collect methods)
  - One collector per sensor, all registered with Prometheus
  - Exposes two metrics:
    - `dht_temperature_degree` (labels: dht_name, hostname, unit)
    - `dht_humidity_percent` (labels: dht_name, hostname)
  - Uses GaugeValue metric type

- **internal/logger/log.go**: Logger factory using logrus
  - Configured via config file with levels: debug, info, warn, error, fatal, panic
  - Uses text formatter with full timestamps, colors disabled

### Data Flow

1. Config read from YAML → Config struct with sensors array
2. Logger initialized with config log level
3. DHT host initialized once via `sensor.HostInit()`
4. For each sensor in config:
   - Sensor initialized with GPIO pin and temperature unit
   - Collector created and registered with Prometheus
5. HTTP server exposes /metrics endpoint
6. Each scrape triggers Collect() on all collectors → sensor.ReadData() → metrics returned

### GPIO Access

- Requires user to be in `gpio` group for GPIO pin access
- GPIO pins specified as integers in config, formatted as "GPIOx" internally
- Uses go-dht library which requires HostInit() before sensor creation

### Systemd Integration

The project includes `examples/dht-prometheus-exporter.service` for systemd management:
```bash
sudo systemctl start dht-prometheus-exporter
sudo systemctl stop dht-prometheus-exporter
sudo systemctl status dht-prometheus-exporter
sudo systemctl enable dht-prometheus-exporter  # Start on boot
```

## Configuration File

Default config file: `dht-prometheus-exporter.yml`

```yaml
sensors:
  - name: living-room      # Sensor identifier (used in metric labels)
    gpio_pin: 2            # GPIO pin number (e.g., 2 for GPIO2)
    max_retries: 10        # Number of read attempts before giving up
    temperature_unit: celsius  # Either "celsius" or "fahrenheit"
  - name: bedroom
    gpio_pin: 17
    max_retries: 10
    temperature_unit: celsius
listen_port: 8080          # HTTP server port
log_level: info            # Logging verbosity (debug, info, warn, error, fatal, panic)
```

## Dependencies

Uses Go modules for dependency management. Key dependencies:
- github.com/prometheus/client_golang
- github.com/MichaelS11/go-dht
- github.com/sirupsen/logrus
- github.com/spf13/viper
