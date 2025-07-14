package logger

type Config struct {
	Level      string `yaml:"level"`       // debug, info, warn, error
	FilePath   string `yaml:"file_path"`   // полный путь к файлу логов (включая имя файла)
	MaxSizeMB  int    `yaml:"max_size"`    // макс. размер в MB
	MaxBackups int    `yaml:"max_backups"` // макс. число файлов
	MaxAgeDays int    `yaml:"max_age"`     // срок хранения в днях
	Compress   bool   `yaml:"compress"`    // сжатие старых логов
}

// DefaultConfig возвращает конфиг по умолчанию
func DefaultConfig() Config {
	return Config{
		Level:      "info",
		FilePath:   "./logs/testlog.log",
		MaxSizeMB:  100,
		MaxBackups: 3,
		MaxAgeDays: 30,
		Compress:   true,
	}
}
