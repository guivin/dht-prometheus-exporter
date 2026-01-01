package main

import (
	"fmt"
	"log"
	stdlibLog "log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/guivin/dht-prometheus-exporter/internal/collector"
	"github.com/guivin/dht-prometheus-exporter/internal/config"
	"github.com/guivin/dht-prometheus-exporter/internal/logger"
	"github.com/guivin/dht-prometheus-exporter/internal/sensor"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	lg, err := logger.New(cfg.LogLevel)
	if err != nil {
		// If logger creation fails, fall back to a default logger
		lg, _ = logger.New("info")
		lg.Warnf("Failed to create logger with level %q, using info level: %v", cfg.LogLevel, err)
	}

	// Initialize sensor
	sensorReader, err := sensor.New(cfg, lg)
	if err != nil {
		return fmt.Errorf("failed to initialize sensor: %w", err)
	}

	// Create collector
	coll := collector.New(sensorReader, cfg.Name, lg)

	// Register collector with Prometheus
	lg.Debug("Registering prometheus collector")
	if err := prometheus.Register(coll); err != nil {
		return fmt.Errorf("failed to register collector: %w", err)
	}

	// Set up HTTP server
	w := lg.Writer()
	defer w.Close()

	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			ErrorLog: stdlibLog.New(w, "", 0),
		},
	))

	// Start server
	addr := fmt.Sprintf(":%d", cfg.ListenPort)
	lg.Infof("Starting HTTP server on %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		return fmt.Errorf("HTTP server error: %w", err)
	}

	return nil
}
