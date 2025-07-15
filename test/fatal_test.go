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
	// Используем временный файл
	testLog := filepath.Join(testLogsDir, "fatal_test.log")

	// Создаем директорию если нужно
	if err := os.MkdirAll(testLogsDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Подменяем os.Exit через замыкание
	var exitCode int
	oldExit := func(code int) { os.Exit(code) }
	defer func() {
		l.SetExitFunc(oldExit) // Восстанавливаем оригинальный exit
	}()

	l.SetExitFunc(func(code int) {
		exitCode = code
		panic("exit called") // Имитируем выход
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

	// Ожидаем вызов os.Exit через panic
	func() {
		defer func() {
			if r := recover(); r != nil && r.(string) != "exit called" {
				panic(r) // Пробрасываем другие panic
			}
		}()
		log.Fatal("fatal error occurred")
	}()

	// Проверяем код выхода
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	// Даем время на запись лога
	time.Sleep(100 * time.Millisecond)

	// Проверяем содержимое лог-файла
	content := readLogFile(t, testLog)
	if !strings.Contains(content, "fatal error occurred") {
		t.Error("Expected fatal message in log")
	}
}

func TestLogger_Panic(t *testing.T) {
	// Используем временный файл
	testLog := filepath.Join(testLogsDir, "panic_test.log")

	cfg := l.Config{
		Level:    "error",
		FilePath: testLog,
	}

	// Проверяем что вызовет panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic, got none")
		}

		// Проверяем что сообщение записалось в лог
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
