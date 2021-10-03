# DHT Prometheus Exporter
> A Prometheus exporter for DHT22/AM2302 sensors runnable on Raspberry Pi

# Prerequisites

Install golang and configure your environment variables:
```
$ sudo apt-get install golang
$ export GOPATH="$HOME/go"
$ export GOBIN="$GOPATH/bin"
$ export PATH="$PATH:$GOBIN"
```

Install make to build the source code:
```
$ sudo apt install make
```

Install dep to manage golang dependencies:
```
$ go get -u github.com/golang/dep/cmd/dep
```

# Installation

On the Raspberry Pi download the project:
```
$ go get -u github.com/guivin/dht-prometheus-exporter
$ cd $GOPATH/src/github.com/guivin/dht-prometheus-exporter
```

Build and install:
```
$ make all
```

Create a dedicated system user and group belonging to gpio group:
```
$ useradd --user-group --groups gpio --no-create-home --system --shell /usr/sbin/nologin dht-prometheus-exporter
```

Copy the configuration file and adapt if needed:
```
$ cp dht-prometheus-exporter.yml /etc/dht-prometheus-exporter.yml
$ sudo chown dht-prometheus-exporter:dht-prometheus-exporter /etc/dht-prometheus-exporter.yml
$ sudo chmod 0640 /etc/dht-prometheus-exporter.yml
```

Configure with systemd:
```
$ cp dht-prometheus-exporter.service /etc/systemd/system
$ sudo systemctl daemon-reload
$ sudo systemctl start dht-prometheus-exporter
```

# Usage

Example with the listenPort configured with TPC/8080:
```
$ curl http://localhost:8080/metrics
```
