package logger_test

import (
	l "logger"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var osExit = os.Exit

func TestLogger_Fatal(t *testing.T) {
	oldExit := osExit
	defer func() { osExit = oldExit }()

	var exitCode int
	osExit = func(code int) { exitCode = code }

	tempDir := t.TempDir()
	cfg := l.Config{
		Level:     "error",
		Directory: tempDir,
		Filename:  "fatal_test.log",
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

	log.Fatal("fatal error occurred")

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	content := readLogFile(t, filepath.Join(tempDir, "fatal_test.log"))
	if !strings.Contains(content, "fatal error occurred") {
		t.Error("Expected fatal message in log")
	}
}

func TestLogger_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic, got none")
		}
	}()

	tempDir := t.TempDir()
	cfg := l.Config{
		Level:     "error",
		Directory: tempDir,
		Filename:  "panic_test.log",
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

	log.Panic("panic situation")

	content := readLogFile(t, filepath.Join(tempDir, "panic_test.log"))
	if !strings.Contains(content, "panic situation") {
		t.Error("Expected panic message in log")
	}
}
