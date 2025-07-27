package logger_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	l "github.com/SerMoskvin/logger"
)

var testDir = func() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}()

func getTestLogPath(filename string) string {
	return filepath.Join(testDir, "test_files", filename)
}

func TestLoadConfig(t *testing.T) {
	t.Run("Should load default config when file not found", func(t *testing.T) {
		_, err := l.LoadConfig("nonexistent_file.yml")
		if err != nil {
			t.Fatalf("Expected to load default config, got error: %v", err)
		}
	})

	t.Run("Should load custom config with correct structure", func(t *testing.T) {
		testConfigPath := filepath.Join("..", "default_config.yml")

		cfg, err := l.LoadConfig(testConfigPath)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if cfg.Debug.FilePath != "../logs/debug.log" {
			t.Errorf("Expected debug file path '../logs/debug.log', got '%s'", cfg.Debug.FilePath)
		}

		if cfg.Info.FilePath != "../logs/info.log" {
			t.Errorf("Expected info file path '../logs/info.log', got '%s'", cfg.Info.FilePath)
		}

		if cfg.Warn.FilePath != "../logs/warn.log" {
			t.Errorf("Expected warn file path '../logs/warn.log', got '%s'", cfg.Warn.FilePath)
		}

		if cfg.Error.FilePath != "../logs/error.log" {
			t.Errorf("Expected error file path '../logs/error.log', got '%s'", cfg.Error.FilePath)
		}
	})
}

func TestNewLevel(t *testing.T) {
	testLogsDir := "test_files"
	testLog := filepath.Join(testLogsDir, "level_logger_test.log")

	if err := os.MkdirAll(testLogsDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.Remove(testLog)

	tmpConfigContent := `
debug:
  level: debug
  file_path: ` + testLog + `
  max_size: 10
  max_backups: 3
  max_age: 7
  compress: false
info:
  level: info
  file_path: ` + testLog + `
  max_size: 10
  max_backups: 3
  max_age: 7
  compress: false
warn:
  level: warn
  file_path: ` + testLog + `
  max_size: 10
  max_backups: 3
  max_age: 7
  compress: false
error:
  level: error
  file_path: ` + testLog + `
  max_size: 10
  max_backups: 3
  max_age: 7
  compress: false
`

	tmpConfigFile := filepath.Join(testLogsDir, "tmp_config.yml")
	if err := os.WriteFile(tmpConfigFile, []byte(tmpConfigContent), 0644); err != nil {
		t.Fatalf("Failed to write temp config file: %v", err)
	}
	defer os.Remove(tmpConfigFile)

	t.Run("Should create level logger with test config file", func(t *testing.T) {
		logger, err := l.NewLevel(tmpConfigFile)
		if err != nil {
			t.Fatalf("Failed to create level logger: %v", err)
		}
		defer logger.Close()

		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")

		if err := logger.Sync(); err != nil {
			t.Fatalf("Failed to sync logs: %v", err)
		}

		content, err := os.ReadFile(testLog)
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}

		messages := []string{
			"debug message",
			"info message",
			"warn message",
			"error message",
		}

		for _, msg := range messages {
			if !strings.Contains(string(content), msg) {
				t.Errorf("Expected to find %q in log", msg)
			}
		}
	})
}

func TestNewLogger(t *testing.T) {
	t.Run("ShouldCreateLoggerWithConfig", func(t *testing.T) {
		var testLogName = "testlog.log"
		cfg := l.Config{
			Level:      "info",
			FilePath:   getTestLogPath(testLogName),
			MaxSizeMB:  10,
			MaxBackups: 3,
			MaxAgeDays: 30,
			Compress:   true,
		}

		if err := os.MkdirAll(filepath.Dir(cfg.FilePath), 0755); err != nil {
			t.Fatalf("Failed to create test_files directory: %v", err)
		}

		_ = os.Remove(cfg.FilePath)

		log, err := l.New(cfg)
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}
		defer cleanupLogger(t, log)

		log.Info("test log creation")

		if err := log.Sync(); err != nil {
			t.Fatalf("Failed to sync logs: %v", err)
		}

		if _, err := os.Stat(cfg.FilePath); os.IsNotExist(err) {
			t.Errorf("Log file was not created at %s", cfg.FilePath)
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
