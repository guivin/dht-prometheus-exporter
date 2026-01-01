package logger

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestParseLevel_ValidLevels(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected log.Level
	}{
		{"debug level", "debug", log.DebugLevel},
		{"info level", "info", log.InfoLevel},
		{"warn level", "warn", log.WarnLevel},
		{"error level", "error", log.ErrorLevel},
		{"fatal level", "fatal", log.FatalLevel},
		{"panic level", "panic", log.PanicLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, err := ParseLevel(tt.input)
			if err != nil {
				t.Errorf("ParseLevel(%q) returned unexpected error: %v", tt.input, err)
			}
			if level != tt.expected {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, level, tt.expected)
			}
		})
	}
}

func TestParseLevel_InvalidLevel(t *testing.T) {
	tests := []string{
		"invalid",
		"DEBUG",  // Case sensitive
		"INFO",   // Case sensitive
		"warning", // Should be "warn"
		"",
		"trace",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			level, err := ParseLevel(input)
			if err == nil {
				t.Errorf("ParseLevel(%q) expected error, but got level: %v", input, level)
			}
		})
	}
}

func TestNew_ValidLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		expected log.Level
	}{
		{"info logger", "info", log.InfoLevel},
		{"debug logger", "debug", log.DebugLevel},
		{"error logger", "error", log.ErrorLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.level)
			if err != nil {
				t.Fatalf("New(%q) returned unexpected error: %v", tt.level, err)
			}
			if logger == nil {
				t.Fatal("New() returned nil logger")
			}
			if logger.Level != tt.expected {
				t.Errorf("Logger level = %v, want %v", logger.Level, tt.expected)
			}
		})
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	logger, err := New("invalid")
	if err == nil {
		t.Error("New(\"invalid\") expected error, got nil")
	}
	if logger != nil {
		t.Errorf("New(\"invalid\") expected nil logger, got %v", logger)
	}
}

func TestNew_OutputConfiguration(t *testing.T) {
	logger, err := New("info")
	if err != nil {
		t.Fatalf("New(\"info\") returned unexpected error: %v", err)
	}

	// Verify logger outputs to stderr
	if logger.Out != os.Stderr {
		t.Errorf("Logger output = %v, want os.Stderr", logger.Out)
	}

	// Verify formatter is TextFormatter
	if _, ok := logger.Formatter.(*log.TextFormatter); !ok {
		t.Errorf("Logger formatter = %T, want *log.TextFormatter", logger.Formatter)
	}

	// Verify formatter configuration
	if formatter, ok := logger.Formatter.(*log.TextFormatter); ok {
		if !formatter.DisableColors {
			t.Error("Expected DisableColors to be true")
		}
		if !formatter.FullTimestamp {
			t.Error("Expected FullTimestamp to be true")
		}
	}
}
