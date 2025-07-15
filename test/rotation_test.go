package logger_test

import (
	l "logger"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestLogger_Rotation(t *testing.T) {
	rotatedLogPath := getTestLogPath("rotation_test.log")

	cfg := l.Config{
		Level:      "info",
		FilePath:   rotatedLogPath,
		MaxSizeMB:  1, // 1 MB
		MaxBackups: 2,
		Compress:   true,
	}

	// Удаляем старые файлы
	removeOldLogFiles(t, rotatedLogPath)

	log, err := l.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer cleanupLogger(t, log)

	// Увеличиваем объём данных (200KB на сообщение, 15 сообщений = 3MB)
	message := strings.Repeat("a", 200000) // 200KB per message
	for i := 0; i < 15; i++ {
		log.Info(message, zap.Int("index", i))
		if i%3 == 0 {
			if err := log.Sync(); err != nil {
				t.Logf("Sync error: %v", err)
			}
		}
	}

	// Увеличиваем время ожидания
	time.Sleep(2 * time.Second)

	// Проверяем результаты
	checkRotationResults(t, rotatedLogPath)
}

func checkRotationResults(t *testing.T, mainLogPath string) {
	dir := filepath.Dir(mainLogPath)
	baseName := strings.TrimSuffix(filepath.Base(mainLogPath), filepath.Ext(mainLogPath))

	// Ищем все файлы, включая сжатые
	files, err := filepath.Glob(filepath.Join(dir, baseName+"*"))
	if err != nil {
		t.Fatalf("Failed to find log files: %v", err)
	}

	// Фильтруем только релевантные файлы
	var logFiles []string
	for _, f := range files {
		if f == mainLogPath ||
			strings.Contains(f, "-") ||
			strings.HasSuffix(f, ".gz") {
			logFiles = append(logFiles, f)
		}
	}

	t.Logf("Found log files:\n%s", strings.Join(logFiles, "\n"))

	// Проверяем размеры файлов
	for _, f := range logFiles {
		info, err := os.Stat(f)
		if err == nil {
			t.Logf("%s: %.2f MB", filepath.Base(f), float64(info.Size())/1024/1024)
		}
	}

	if len(logFiles) < 2 {
		t.Errorf("Expected at least 2 files (current + rotated), got %d", len(logFiles))
	}
}

// Вспомогательная функция для удаления старых лог-файлов
func removeOldLogFiles(t *testing.T, basePath string) {
	dir := filepath.Dir(basePath)
	base := filepath.Base(basePath)
	baseName := strings.TrimSuffix(base, filepath.Ext(base))

	files, err := filepath.Glob(filepath.Join(dir, baseName+"*"))
	if err != nil {
		t.Logf("Warning: failed to clean old log files: %v", err)
		return
	}

	for _, f := range files {
		if err := os.Remove(f); err != nil {
			t.Logf("Warning: failed to remove old log file %s: %v", f, err)
		}
	}
}

// Вспомогательная функция для очистки логгера
func cleanupLogger(t *testing.T, log *l.Logger) {
	if err := log.Sync(); err != nil {
		t.Logf("Sync error: %v", err)
	}
	if err := log.Close(); err != nil {
		t.Logf("Close error: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
}
