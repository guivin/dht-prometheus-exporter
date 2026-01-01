# Example Configuration Files

This directory contains example configuration files for the DHT Prometheus Exporter.

## Files

### dht-prometheus-exporter.yml

Example configuration file for the exporter. Copy this to `/etc/dht-prometheus-exporter.yml` and modify according to your setup.

**Configuration options:**

- `name`: Sensor name used in Prometheus metrics labels
- `gpio_pin`: GPIO pin number where the DHT22/AM2302 sensor is connected (e.g., 2, 4, 17)
- `max_retries`: Number of retry attempts when reading from the sensor (recommended: 10)
- `listen_port`: HTTP port for the metrics endpoint (default: 8080)
- `log_level`: Logging verbosity - one of: debug, info, warn, error, fatal, panic
- `temperature_unit`: Temperature unit - either `celsius` or `fahrenheit`

**Example usage:**

```bash
sudo cp examples/dht-prometheus-exporter.yml /etc/dht-prometheus-exporter.yml
sudo chown dht-prometheus-exporter:dht-prometheus-exporter /etc/dht-prometheus-exporter.yml
sudo chmod 0640 /etc/dht-prometheus-exporter.yml
```

### dht-prometheus-exporter.service

Example systemd service file for running the exporter as a system service.

The service is configured to:
- Run as the `dht-prometheus-exporter` user with `gpio` group access
- Automatically restart on failure
- Log to syslog with identifier `dht-prometheus-exporter`
- Start after network and filesystem are available

**Example usage:**

```bash
sudo cp examples/dht-prometheus-exporter.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable dht-prometheus-exporter
sudo systemctl start dht-prometheus-exporter
```

Check logs:

```bash
sudo journalctl -u dht-prometheus-exporter -f
```
