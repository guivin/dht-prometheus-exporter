# DHT Prometheus Exporter
> A Prometheus exporter for DHT22/AM2302 sensors runnable on Raspberry Pi

This repository contains a Prometheus exporter designed for the DHT22/AM2302 temperature and humidity sensors, optimized 
for use on Raspberry Pi devices.

## Prerequisites

Before you begin, ensure you have the necessary tools and dependencies installed:

* Install Golang and set up your environment variables:

```
sudo apt-get install golang
export GOPATH="$HOME/go"
export GOBIN="$GOPATH/bin"
export PATH="$PATH:$GOBIN"
```

* Make: Required for building the source code:

```
sudo apt install make
```

* Dep: Golang's dependency management tool:

```
go get -u github.com/golang/dep/cmd/dep
```

## Installation

Follow these steps to install the DHT Prometheus Exporter on your Raspberry Pi:

1. Download the project:

```
go get -u github.com/guivin/dht-prometheus-exporter
cd $GOPATH/src/github.com/guivin/dht-prometheus-exporter
```

2. Build and install the project:

```
make all
```

3. Create a dedicated system user that belongs to the gpio group (for GPIO pin access):

```
useradd --user-group --groups gpio --no-create-home --system --shell /usr/sbin/nologin dht-prometheus-exporter
```

4. Set up the configuration file. Copy the default configuration file and modify it according to your needs:

```
cp dht-prometheus-exporter.yml /etc/dht-prometheus-exporter.yml
sudo chown dht-prometheus-exporter:dht-prometheus-exporter /etc/dht-prometheus-exporter.yml
sudo chmod 0640 /etc/dht-prometheus-exporter.yml
```

5. Integrate with systemd for easy service management:

```
cp dht-prometheus-exporter.service /etc/systemd/system
sudo systemctl daemon-reload
sudo systemctl start dht-prometheus-exporter
```

## Usage

Retrieve the metrics from the exporter by querying the designated HTTP endpoint (adjust the port if 
your configuration differs):

```http
  GET http://localhost:8080/metrics
```

This command will output the current readings from your DHT22/AM2302 sensors, making the data available for Prometheus 
scraping and subsequent analysis or visualization.
