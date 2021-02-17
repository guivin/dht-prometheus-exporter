package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

const defaultLogLevel = log.InfoLevel

var (
	once     sync.Once
	instance log.Logger
)

func LogLevel(levelName string) log.Level {
	/**
	Set the log level of logrus following the value supplied in the configuration
	*/
	var logLevel log.Level

	switch levelName {
	case "debug":
		logLevel = log.DebugLevel
	case "info":
		logLevel = log.InfoLevel
	case "warn":
		logLevel = log.WarnLevel
	case "error":
		logLevel = log.ErrorLevel
	case "fatal":
		logLevel = log.FatalLevel
	case "panic":
		logLevel = log.PanicLevel
	default:
		logLevel = defaultLogLevel
	}

	return logLevel
}

func getLogger(config *Config) log.Logger {
	/**
	Singleton to get the same logger instance everywhere
	*/
	once.Do(func() {
		instance = log.Logger{
			Out:   os.Stderr,
			Level: LogLevel(config.defaultLogLevel),
			Formatter: &log.TextFormatter{
				DisableColors: true,
				FullTimestamp: true,
			},
		}
	})

	return instance
}
