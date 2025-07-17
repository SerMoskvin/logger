package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type LevelLogger struct {
	debug *Logger
	info  *Logger
	warn  *Logger
	error *Logger
	mu    sync.Mutex
}

func NewLevel(cfg LevelConfig) (*LevelLogger, error) {
	debugLog, err := New(cfg.Debug)
	if err != nil {
		return nil, fmt.Errorf("failed to create debug logger: %w", err)
	}

	infoLog, err := New(cfg.Info)
	if err != nil {
		return nil, fmt.Errorf("failed to create info logger: %w", err)
	}

	warnLog, err := New(cfg.Warn)
	if err != nil {
		return nil, fmt.Errorf("failed to create warn logger: %w", err)
	}

	errorLog, err := New(cfg.Error)
	if err != nil {
		return nil, fmt.Errorf("failed to create error logger: %w", err)
	}

	return &LevelLogger{
		debug: debugLog,
		info:  infoLog,
		warn:  warnLog,
		error: errorLog,
	}, nil
}

func (m *LevelLogger) Debug(msg string, fields ...zap.Field) {
	m.debug.Debug(msg, fields...)
}

func (m *LevelLogger) Info(msg string, fields ...zap.Field) {
	m.info.Info(msg, fields...)
}

func (m *LevelLogger) Warn(msg string, fields ...zap.Field) {
	m.warn.Warn(msg, fields...)
}

func (m *LevelLogger) Error(msg string, fields ...zap.Field) {
	m.error.Error(msg, fields...)
}

func (m *LevelLogger) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	if err := m.debug.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := m.info.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := m.warn.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := m.error.Close(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing loggers: %v", errs)
	}
	return nil
}

func (m *LevelLogger) Sync() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	if err := m.debug.Sync(); err != nil {
		errs = append(errs, err)
	}
	if err := m.info.Sync(); err != nil {
		errs = append(errs, err)
	}
	if err := m.warn.Sync(); err != nil {
		errs = append(errs, err)
	}
	if err := m.error.Sync(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors syncing loggers: %v", errs)
	}
	return nil
}
