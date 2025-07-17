package logger

type Config struct {
	Level      string `yaml:"level"`
	FilePath   string `yaml:"file_path"`
	MaxSizeMB  int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAgeDays int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

type LevelConfig struct {
	Debug Config `yaml:"debug"`
	Info  Config `yaml:"info"`
	Warn  Config `yaml:"warn"`
	Error Config `yaml:"error"`
}

func DefaultConfig() Config {
	return Config{
		Level:      "info",
		FilePath:   "./logs/app.log",
		MaxSizeMB:  10,
		MaxBackups: 3,
		MaxAgeDays: 30,
		Compress:   true,
	}
}

func DefaultLevelConfig() LevelConfig {
	return LevelConfig{
		Debug: Config{
			Level:      "debug",
			FilePath:   "./logs/debug.log",
			MaxSizeMB:  10,
			MaxBackups: 3,
			MaxAgeDays: 7,
			Compress:   false,
		},
		Info: Config{
			Level:      "info",
			FilePath:   "./logs/info.log",
			MaxSizeMB:  20,
			MaxBackups: 5,
			MaxAgeDays: 30,
			Compress:   true,
		},
		Warn: Config{
			Level:      "warn",
			FilePath:   "./logs/warn.log",
			MaxSizeMB:  30,
			MaxBackups: 7,
			MaxAgeDays: 60,
			Compress:   true,
		},
		Error: Config{
			Level:      "error",
			FilePath:   "./logs/error.log",
			MaxSizeMB:  50,
			MaxBackups: 10,
			MaxAgeDays: 90,
			Compress:   true,
		},
	}
}
