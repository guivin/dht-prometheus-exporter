# DHT Prometheus Exporter

[![Go Version](https://img.shields.io/github/go-mod/go-version/guivin/dht-prometheus-exporter)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/guivin/dht-prometheus-exporter)](https://goreportcard.com/report/github.com/guivin/dht-prometheus-exporter)
[![License](https://img.shields.io/github/license/guivin/dht-prometheus-exporter)](LICENSE)
[![GitHub release](https://img.shields.io/github/release/guivin/dht-prometheus-exporter.svg)](https://github.com/guivin/dht-prometheus-exporter/releases)
[![GitHub issues](https://img.shields.io/github/issues/guivin/dht-prometheus-exporter)](https://github.com/guivin/dht-prometheus-exporter/issues)

<p align="center">
  <img src="assets/raspberry-dht22.jpeg" alt="Raspberry Pi with DHT22 Sensor" height="200">
</p>

> A Prometheus exporter for DHT22/AM2302 temperature and humidity sensors designed for Raspberry Pi

This repository contains a production-ready Prometheus exporter for the DHT22/AM2302 temperature and humidity sensors, optimized
for use on Raspberry Pi devices. Built with Go, it provides reliable metrics collection with proper error handling and retry logic.

## Prerequisites

Before you begin, ensure you have the necessary tools and dependencies installed:

* **Go 1.16 or later**: Required for building the project with Go modules:

```bash
sudo apt-get install golang
```

* **Make**: Required for building the source code:

```bash
sudo apt install make
```

## Installation

Follow these steps to install the DHT Prometheus Exporter on your Raspberry Pi:

1. Clone the repository:

```bash
git clone https://github.com/guivin/dht-prometheus-exporter.git
cd dht-prometheus-exporter
```

2. Build the project:

```bash
make build
```

Or build manually:

```bash
go build -o dht-prometheus-exporter ./cmd/dht-prometheus-exporter
```

3. Install the binary (optional):

```bash
make install
```

4. Create a dedicated system user that belongs to the gpio group (for GPIO pin access):

```
useradd --user-group --groups gpio --no-create-home --system --shell /usr/sbin/nologin dht-prometheus-exporter
```

5. Set up the configuration file. Copy the example configuration file and modify it according to your needs:

```bash
sudo cp examples/dht-prometheus-exporter.yml /etc/dht-prometheus-exporter.yml
sudo chown dht-prometheus-exporter:dht-prometheus-exporter /etc/dht-prometheus-exporter.yml
sudo chmod 0640 /etc/dht-prometheus-exporter.yml
```

Edit `/etc/dht-prometheus-exporter.yml` to configure:
- `name`: Sensor name for metrics labels
- `gpio_pin`: GPIO pin number where DHT22 is connected
- `max_retries`: Number of retry attempts for sensor reads
- `listen_port`: HTTP port for metrics endpoint (default: 8080)
- `log_level`: Logging level (debug, info, warn, error)
- `temperature_unit`: celsius or fahrenheit

6. Integrate with systemd for easy service management:

```bash
sudo cp examples/dht-prometheus-exporter.service /etc/systemd/system
sudo systemctl daemon-reload
sudo systemctl enable dht-prometheus-exporter
sudo systemctl start dht-prometheus-exporter
```

Check service status:

```bash
sudo systemctl status dht-prometheus-exporter
```

## Usage

Retrieve the metrics from the exporter by querying the designated HTTP endpoint (adjust the port if
your configuration differs):

```http
GET http://localhost:8080/metrics
```

This command will output the current readings from your DHT22/AM2302 sensors, making the data available for Prometheus
scraping and subsequent analysis or visualization.

## Testing

Run the test suite:

```bash
make test
```

Generate test coverage report:

```bash
make test-coverage
```

This will create `coverage.html` that you can open in a browser to see detailed coverage information.

## Development

### Project Structure

The project follows standard Go conventions:

```
dht-prometheus-exporter/
├── cmd/
│   └── dht-prometheus-exporter/    # Application entry point
├── internal/                        # Internal packages
│   ├── config/                      # Configuration management
│   ├── sensor/                      # DHT sensor interface and implementation
│   ├── collector/                   # Prometheus collector
│   └── logger/                      # Logging configuration
├── examples/                        # Example configuration files
│   ├── dht-prometheus-exporter.yml # Example config file
│   └── dht-prometheus-exporter.service # Example systemd service
├── testdata/                        # Test data files
├── go.mod                           # Go module definition
└── Makefile                         # Build automation
```

### Building from Source

```bash
# Build the binary
make build

# Run tests
make test

# Clean build artifacts
make clean

# Tidy Go modules
make mod-tidy
```
