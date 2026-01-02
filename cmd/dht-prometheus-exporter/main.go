package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	stdlibLog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Initialize DHT host (required before creating sensors)
	lg.Info("Initializing DHT22/AM2302 host")
	if err := sensor.HostInit(); err != nil {
		return fmt.Errorf("failed to initialize DHT host: %w", err)
	}

	// Initialize sensors and collectors
	for i := range cfg.Sensors {
		sensorCfg := &cfg.Sensors[i]
		sensorReader, err := sensor.New(sensorCfg, lg)
		if err != nil {
			return fmt.Errorf("failed to initialize sensor '%s': %w", sensorCfg.Name, err)
		}

		// Create and register collector
		coll := collector.New(sensorReader, lg)
		lg.Debugf("Registering prometheus collector for sensor '%s'", sensorCfg.Name)
		if err := prometheus.Register(coll); err != nil {
			return fmt.Errorf("failed to register collector for sensor '%s': %w", sensorCfg.Name, err)
		}
	}

	lg.Infof("Initialized %d sensor(s)", len(cfg.Sensors))

	// Set up HTTP server
	w := lg.Writer()
	defer func() { _ = w.Close() }()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			ErrorLog: stdlibLog.New(w, "", 0),
		},
	))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := fmt.Sprintf(":%d", cfg.ListenPort)
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for shutdown signals
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		lg.Info("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			lg.Errorf("Could not gracefully shutdown the server: %v", err)
		}
		close(done)
	}()

	lg.Infof("Starting HTTP server on %s", addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server error: %w", err)
	}

	<-done
	lg.Info("Server stopped")

	return nil
}
