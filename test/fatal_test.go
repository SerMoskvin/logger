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
	testLogsDir := "test_files"
	testLog := filepath.Join(testLogsDir, "fatal_test.log")

	if err := os.MkdirAll(testLogsDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testLogsDir)

	originalExit := func(code int) { os.Exit(code) }
	defer l.SetExitFunc(originalExit)

	exitCalled := false
	exitCode := 0

	l.SetExitFunc(func(code int) {
		exitCalled = true
		exitCode = code
		panic("fatal exit")
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
		if err := log.Close(); err != nil {
			t.Logf("Close error: %v", err)
		}
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				if r != "fatal exit" {
					t.Errorf("Unexpected panic: %v", r)
				}
			}
		}()
		log.Fatal("fatal error occurred")
	}()

	if !exitCalled {
		t.Error("Expected exit function to be called")
	}

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	time.Sleep(100 * time.Millisecond)

	content, err := os.ReadFile(testLog)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "fatal error occurred") {
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
