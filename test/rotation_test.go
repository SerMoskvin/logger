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
	rotatedLogPath := `C:\Users\Joker\Desktop\logger\test\test_files\rotation_test.log`

	cfg := l.Config{
		Level:      "info",
		FilePath:   rotatedLogPath,
		MaxSizeMB:  1, // 1 MB
		MaxBackups: 2,
		Compress:   true,
	}

	log, err := l.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer cleanupLogger(t, log)

	message := strings.Repeat("a", 100000) // 100KB per message
	for i := 0; i < 20; i++ {
		log.Info(message, zap.Int("index", i))
		// Принудительная синхронизация каждые 5 сообщений
		if i%5 == 0 {
			if err := log.Sync(); err != nil {
				t.Logf("Sync error: %v", err)
			}
		}
	}

	time.Sleep(1 * time.Second)

	dir := filepath.Dir(rotatedLogPath)
	base := filepath.Base(rotatedLogPath)
	baseName := strings.TrimSuffix(base, filepath.Ext(base))

	// Ищем все файлы, начинающиеся с базового имени
	allFiles, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	var matchedFiles []string
	for _, file := range allFiles {
		fileName := file.Name()
		if strings.HasPrefix(fileName, baseName) {
			matchedFiles = append(matchedFiles, filepath.Join(dir, fileName))
		}
	}

	t.Logf("All matched files: %v", matchedFiles)

	// Ожидаем как минимум 2 файла
	if len(matchedFiles) < 2 {
		fileInfo, err := os.Stat(rotatedLogPath)
		if err == nil {
			t.Logf("Main log file size: %d bytes", fileInfo.Size())
		}
		t.Errorf("Expected at least 2 files (current + rotated), got %d", len(matchedFiles))
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
