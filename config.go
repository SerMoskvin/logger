package logger

import (
	"embed"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

//go:embed default_config.yml
var defaultConfigFS embed.FS

const defaultConfigPath = "default_config.yml"

func LoadDefaultConfig() (*LevelConfig, error) {
	return LoadConfig("")
}

// Field представляет собой лог-поле
type Field struct {
	Key   string
	Value interface{}
}

// Config содержит конфигурацию логгера
type Config struct {
	Level      string `yaml:"level"`
	FilePath   string `yaml:"file_path"`
	MaxSizeMB  int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAgeDays int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

// LevelConfig содержит конфигурацию для разных уровней логирования
type LevelConfig struct {
	Debug Config `yaml:"debug"`
	Info  Config `yaml:"info"`
	Warn  Config `yaml:"warn"`
	Error Config `yaml:"error"`
}

// LoadConfig загружает конфигурацию из YAML файла
func LoadConfig(configPath string) (*LevelConfig, error) {
	var data []byte
	var err error

	data, err = os.ReadFile(configPath)
	if err != nil {
		data, err = defaultConfigFS.ReadFile(defaultConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read both custom and default config: %w", err)
		}
	}

	var cfg LevelConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

var exitFunc = os.Exit

func SetExitFunc(f func(int)) {
	exitFunc = f
}
