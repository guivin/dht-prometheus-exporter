# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Prometheus exporter for DHT22/AM2302 temperature and humidity sensors, designed to run on Raspberry Pi. It exposes sensor readings via HTTP endpoint for Prometheus scraping.

## Build and Development Commands

### Build
```bash
make build           # Build the binary
make all            # Full build: dependencies + build + install + cleanup
```

### Install/Uninstall
```bash
make install        # Copy binary to /usr/bin (requires sudo)
make uninstall      # Remove binary from /usr/bin (requires sudo)
```

### Dependencies
```bash
make dep            # Ensure dependencies with dep tool
```

### Cleanup
```bash
make clean          # Remove build artifacts
```

## Architecture

### Core Components

- **main.go**: Entry point that orchestrates initialization and HTTP server setup
  - Reads config from YAML (searches /etc, $HOME, and current directory)
  - Initializes logger, sensor, and collector
  - Starts HTTP server on configured port with /metrics endpoint

- **config.go**: Configuration management using Viper
  - Supports YAML config file named `dht-prometheus-exporter.yml`
  - Config struct contains: name, gpio_pin, max_retries, listen_port, log_level, temperature_unit

- **sensor.go**: DHT22/AM2302 sensor interface wrapper
  - Uses go-dht library for GPIO communication
  - Supports both Celsius and Fahrenheit
  - Provides retry mechanism for reliable sensor reads
  - Requires host initialization before creating sensor instances

- **collector.go**: Prometheus collector implementation
  - Implements prometheus.Collector interface (Describe/Collect methods)
  - Exposes two metrics:
    - `dht_temperature_degree` (labels: dht_name, hostname, unit)
    - `dht_humidity_percent` (labels: dht_name, hostname)
  - Uses CounterValue metric type (note: these are observations, not actual counters)

- **log.go**: Singleton logger using logrus
  - Configured via config file with levels: debug, info, warn, error, fatal, panic
  - Uses text formatter with full timestamps, colors disabled

### Data Flow

1. Config read from YAML → Config struct
2. Logger initialized with config log level
3. Sensor initialized with GPIO pin and temperature unit
4. Collector wraps sensor for Prometheus metrics
5. HTTP server exposes /metrics endpoint
6. Each scrape triggers Collect() → sensor.readRetry() → metrics returned

### GPIO Access

- Requires user to be in `gpio` group for GPIO pin access
- GPIO pins specified as integers in config, formatted as "GPIOx" internally
- Uses go-dht library which requires HostInit() before sensor creation

### Systemd Integration

The project includes `dht-prometheus-exporter.service` for systemd management. Standard systemd commands:
```bash
sudo systemctl start dht-prometheus-exporter
sudo systemctl stop dht-prometheus-exporter
sudo systemctl status dht-prometheus-exporter
sudo systemctl enable dht-prometheus-exporter  # Start on boot
```

## Configuration File

Default config file: `dht-prometheus-exporter.yml`

Required fields:
- `name`: Sensor identifier (used in metric labels)
- `gpio_pin`: GPIO pin number (e.g., 2 for GPIO2)
- `max_retries`: Number of read attempts before giving up
- `listen_port`: HTTP server port (typically 8080)
- `log_level`: Logging verbosity (debug, info, warn, error, fatal, panic)
- `temperature_unit`: Either "celsius" or "fahrenheit"

## Dependencies

Uses `dep` for dependency management (legacy tool, predates Go modules). Key dependencies:
- github.com/prometheus/client_golang
- github.com/MichaelS11/go-dht
- github.com/sirupsen/logrus
- github.com/spf13/viper
