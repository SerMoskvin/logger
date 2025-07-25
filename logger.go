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

// Logger обертка над zap.Logger
type Logger struct {
	*zap.Logger
	level  string
	closer io.Closer
}

// New создает новый экземпляр логгера
func New(cfg Config) (*Logger, error) {
	dir := filepath.Dir(cfg.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	zapLevel := zapcore.InfoLevel
	switch cfg.Level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	}

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

// Методы Logger
func (l *Logger) Debug(msg string, fields ...Field) {
	l.Logger.Debug(msg, convertFields(fields)...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.Logger.Info(msg, convertFields(fields)...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.Logger.Warn(msg, convertFields(fields)...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.Logger.Error(msg, convertFields(fields)...)
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

func (l *Logger) Sync() error {
	return l.Logger.Sync()
}
