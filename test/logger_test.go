package logger_test

import (
	"encoding/json"
	l "logger"
	"os"
	"strings"
	"testing"

	"go.uber.org/zap"
)

const testLogPath = `C:\Users\Joker\Desktop\logger\test\test_files\testlog.log`

func TestNewLogger(t *testing.T) {
	t.Run("ShouldCreateLoggerWithDefaultConfig", func(t *testing.T) {
		cfg := l.DefaultConfig()
		cfg.FilePath = testLogPath

		_ = os.Remove(testLogPath)

		log, err := l.New(cfg)
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}
		defer cleanupLogger(t, log)

		// Делаем тестовую запись в лог
		log.Info("test log creation")

		if err := log.Sync(); err != nil {
			t.Fatalf("Failed to sync logs: %v", err)
		}

		if _, err := os.Stat(testLogPath); os.IsNotExist(err) {
			t.Errorf("Log file was not created at %s", testLogPath)
		}
	})

	t.Run("ShouldErrorWhenCannotCreateDir", func(t *testing.T) {
		cfg := l.Config{
			FilePath: `F:/invalidx22\path\test.log`,
		}
		_, err := l.New(cfg)
		if err == nil {
			t.Error("Expected error when directory cannot be created")
		}
	})
}

func TestLogger_Levels(t *testing.T) {
	// Для этого теста используем отдельный файл
	testLog := `C:\Users\Joker\Desktop\logger\test\test_files\levels_test.log`
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
				log.Debug(tt.message)
			case "info":
				log.Info(tt.message)
			case "warn":
				log.Warn(tt.message)
			case "error":
				log.Error(tt.message)
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
	_ = os.Remove(testLogPath)

	cfg := l.Config{
		Level:    "debug",
		FilePath: testLogPath,
	}

	log, err := l.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer cleanupLogger(t, log)

	testMsg := "test log message"
	log.Info(testMsg, zap.String("key", "value"))

	if err := log.Sync(); err != nil {
		t.Fatalf("Failed to sync logs: %v", err)
	}

	content := readLogFile(t, testLogPath)
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

// Вспомогательная функция для чтения лог-файла
func readLogFile(t *testing.T, path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	return string(data)
}
