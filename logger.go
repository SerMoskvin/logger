package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*zap.Logger
	level  string
	closer io.Closer
}

func New(cfg Config) (*Logger, error) {
	dir := filepath.Dir(cfg.FilePath)

	// Создаем все необходимые директории
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Настройка уровней
	zapLevel := zapcore.InfoLevel
	switch cfg.Level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	}

	// Настройка ротации (используем полный путь из конфига)
	lumberjackLogger := &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   cfg.Compress,
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(lumberjackLogger),
		zapLevel,
	)

	return &Logger{
		Logger: zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)),
		level:  cfg.Level,
		closer: lumberjackLogger,
	}, nil
}

func (l *Logger) Close() error {
	if err := l.Sync(); err != nil {
		return err
	}
	if l.closer != nil {
		return l.closer.Close()
	}
	return nil
}

// Sync синхронизирует буферы
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

var exitFunc = os.Exit

func SetExitFunc(f func(int)) {
	exitFunc = f
}
