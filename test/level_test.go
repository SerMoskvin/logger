package logger_test

import (
	"encoding/json"
	"fmt"
	l "logger"
	"os"
	"strings"
	"testing"
)

func TestLogger_Levels(t *testing.T) {
	testLog := getTestLogPath("levels_test.log")
	_ = os.Remove(testLog)

	cfg := l.Config{
		Level:    "debug",
		FilePath: testLog,
	}

	log, err := l.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer cleanupLogger(t, log)

	tests := []struct {
		level     string
		message   string
		shouldLog bool
	}{
		{"debug", "debug message", true},
		{"info", "info message", true},
		{"warn", "warn message", true},
		{"error", "error message", true},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			switch tt.level {
			case "debug":
				log.Debug(tt.message, l.String("test", "value"))
			case "info":
				log.Info(tt.message, l.Int("test", 123))
			case "warn":
				log.Warn(tt.message, l.Bool("test", true))
			case "error":
				log.Error(tt.message, l.Error(fmt.Errorf("test error")))
			}

			if err := log.Sync(); err != nil {
				t.Fatalf("Failed to sync logs: %v", err)
			}

			content := readLogFile(t, testLog)
			if tt.shouldLog && !strings.Contains(content, tt.message) {
				t.Errorf("Expected to find %q in log", tt.message)
			}
		})
	}
}

func TestLogger_LogFormat(t *testing.T) {
	testLog := getTestLogPath("format_test.log")
	_ = os.Remove(testLog)

	cfg := l.Config{
		Level:    "debug",
		FilePath: testLog,
	}

	log, err := l.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer cleanupLogger(t, log)

	testMsg := "test log message"
	log.Info(testMsg, l.String("key", "value"))

	if err := log.Sync(); err != nil {
		t.Fatalf("Failed to sync logs: %v", err)
	}

	content := readLogFile(t, testLog)
	lines := strings.Split(strings.TrimSpace(content), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		var entry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatalf("Invalid JSON format: %v\n%s", err, line)
		}

		requiredFields := []string{"level", "timestamp", "message"}
		for _, field := range requiredFields {
			if _, ok := entry[field]; !ok {
				t.Errorf("Missing required field: %s", field)
			}
		}

		if entry["message"] != testMsg {
			t.Errorf("Expected message %q, got %q", testMsg, entry["message"])
		}

		if fields, ok := entry["key"].(string); !ok || fields != "value" {
			t.Errorf("Expected field 'key' with value 'value', got %v", entry["key"])
		}
	}
}
