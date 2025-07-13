package logger_test

import (
	"encoding/json"
	l "logger"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	t.Run("ShouldCreateLoggerWithDefaultConfig", func(t *testing.T) {
		cfg := l.DefaultConfig()
		log, err := l.New(cfg)
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}
		defer log.Sync()

		if log == nil {
			t.Error("Expected logger instance, got nil")
		}
	})

	t.Run("ShouldErrorWhenCannotCreateDir", func(t *testing.T) {
		cfg := l.Config{
			Directory: "\x00invalid_path/invalid_path",
			Filename:  "test.log",
		}
		_, err := l.New(cfg)
		if err == nil {
			t.Error("Expected error when directory cannot be created")
		}
	})
}

func TestLogger_Levels(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		level    string
		messages []struct {
			level     string
			message   string
			shouldLog bool
		}
	}{
		{
			name:  "DebugLevel_LogsEverything",
			level: "debug",
			messages: []struct {
				level     string
				message   string
				shouldLog bool
			}{
				{"debug", "debug message", true},
				{"info", "info message", true},
				{"warn", "warn message", true},
				{"error", "error message", true},
			},
		},
		{
			name:  "InfoLevel_FiltersDebug",
			level: "info",
			messages: []struct {
				level     string
				message   string
				shouldLog bool
			}{
				{"debug", "debug message", false},
				{"info", "info message", true},
				{"warn", "warn message", true},
				{"error", "error message", true},
			},
		},
		{
			name:  "ErrorLevel_OnlyErrors",
			level: "error",
			messages: []struct {
				level     string
				message   string
				shouldLog bool
			}{
				{"debug", "debug message", false},
				{"info", "info message", false},
				{"warn", "warn message", false},
				{"error", "error message", true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := l.Config{
				Level:     tt.level,
				Directory: tempDir,
				Filename:  tt.name + ".log",
			}

			log, err := l.New(cfg)
			if err != nil {
				t.Fatalf("Failed to create logger: %v", err)
			}

			defer func() {
				if err := log.Sync(); err != nil {
					t.Logf("Sync error: %v", err)
				}
				if err := log.Close(); err != nil {
					t.Logf("Close error: %v", err)
				}
				time.Sleep(100 * time.Millisecond)
			}()

			for _, msg := range tt.messages {
				switch msg.level {
				case "debug":
					log.Debug(msg.message)
				case "info":
					log.Info(msg.message)
				case "warn":
					log.Warn(msg.message)
				case "error":
					log.Error(msg.message)
				}
			}

			if err := log.Sync(); err != nil {
				t.Fatalf("Failed to sync logs: %v", err)
			}

			content := readLogFile(t, filepath.Join(tempDir, tt.name+".log"))

			for _, msg := range tt.messages {
				if msg.shouldLog {
					if !strings.Contains(content, msg.message) {
						t.Errorf("Expected to find %q in log", msg.message)
					}
				} else {
					if strings.Contains(content, msg.message) {
						t.Errorf("Unexpected %q in log", msg.message)
					}
				}
			}
		})
	}
}

func TestLogger_LogFormat(t *testing.T) {
	tempDir := t.TempDir()
	cfg := l.Config{
		Level:     "debug",
		Directory: tempDir,
		Filename:  "format_test.log",
	}

	log, err := l.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer func() {
		if err := log.Sync(); err != nil {
			t.Logf("Sync error: %v", err)
		}
	}()

	testMsg := "test log message"
	log.Info(testMsg, zap.String("key", "value"))

	if err := log.Sync(); err != nil {
		t.Fatalf("Failed to sync logs: %v", err)
	}

	content := readLogFile(t, filepath.Join(tempDir, "format_test.log"))
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

func TestLogger_Rotation(t *testing.T) {
	tempDir := t.TempDir()
	cfg := l.Config{
		Level:      "info",
		Directory:  tempDir,
		Filename:   "rotation_test.log",
		MaxSizeMB:  1,
		MaxBackups: 2,
		Compress:   true,
	}

	log, err := l.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer func() {
		if err := log.Sync(); err != nil {
			t.Logf("Sync error: %v", err)
		}
	}()

	for i := 0; i < 1000; i++ {
		log.Info(strings.Repeat("test log message ", 50), zap.Int("index", i))
	}

	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read log dir: %v", err)
	}

	var rotatedFiles int
	for _, file := range files {
		if file.Name() != "rotation_test.log" && strings.HasPrefix(file.Name(), "rotation_test.log") {
			rotatedFiles++
		}
	}

	if rotatedFiles != cfg.MaxBackups {
		t.Errorf("Expected %d rotated files, got %d", cfg.MaxBackups, rotatedFiles)
	}
}

func readLogFile(t *testing.T, path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	return string(data)
}
