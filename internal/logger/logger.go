package logger

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// New creates a new logger with the specified log level.
// Returns an error if the log level is invalid.
func New(level string) (*log.Logger, error) {
	logLevel, err := ParseLevel(level)
	if err != nil {
		return nil, err
	}

	logger := log.New()
	logger.SetOutput(os.Stderr)
	logger.SetLevel(logLevel)
	logger.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	return logger, nil
}

// ParseLevel converts a string log level to logrus.Level.
// Valid levels are: debug, info, warn, error, fatal, panic.
// Returns an error for invalid levels.
func ParseLevel(levelName string) (log.Level, error) {
	switch levelName {
	case "debug":
		return log.DebugLevel, nil
	case "info":
		return log.InfoLevel, nil
	case "warn":
		return log.WarnLevel, nil
	case "error":
		return log.ErrorLevel, nil
	case "fatal":
		return log.FatalLevel, nil
	case "panic":
		return log.PanicLevel, nil
	default:
		return 0, fmt.Errorf("invalid log level: %s", levelName)
	}
}
