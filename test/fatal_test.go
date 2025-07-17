package logger_test

import (
	l "logger"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var testLogsDir = filepath.Join("test_files")

func TestLogger_Fatal(t *testing.T) {
	testLog := filepath.Join(testLogsDir, "fatal_test.log")

	if err := os.MkdirAll(testLogsDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	var exitCode int
	oldExit := func(code int) { os.Exit(code) }
	defer func() {
		l.SetExitFunc(oldExit)
	}()

	l.SetExitFunc(func(code int) {
		exitCode = code
		panic("exit called")
	})

	cfg := l.Config{
		Level:    "error",
		FilePath: testLog,
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

	func() {
		defer func() {
			if r := recover(); r != nil && r.(string) != "exit called" {
				panic(r) // Пробрасываем другие panic
			}
		}()
		log.Fatal("fatal error occurred")
	}()

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	time.Sleep(100 * time.Millisecond)

	content := readLogFile(t, testLog)
	if !strings.Contains(content, "fatal error occurred") {
		t.Error("Expected fatal message in log")
	}
}

func TestLogger_Panic(t *testing.T) {
	testLog := filepath.Join(testLogsDir, "panic_test.log")

	cfg := l.Config{
		Level:    "error",
		FilePath: testLog,
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic, got none")
		}

		time.Sleep(100 * time.Millisecond)
		content := readLogFile(t, testLog)
		if !strings.Contains(content, "panic situation") {
			t.Error("Expected panic message in log")
		}
	}()

	log, err := l.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer func() {
		if err := log.Sync(); err != nil {
			t.Logf("Sync error: %v", err)
		}
	}()

	log.Panic("panic situation")
}
