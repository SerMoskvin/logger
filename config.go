package logger

type Config struct {
	Level      string `yaml:"level"`       // debug, info, warn, error
	Directory  string `yaml:"directory"`   // путь к папке с логами
	Filename   string `yaml:"filename"`    // имя файла (app.log)
	MaxSizeMB  int    `yaml:"max_size"`    // макс. размер в MB
	MaxBackups int    `yaml:"max_backups"` // макс. число файлов
	MaxAgeDays int    `yaml:"max_age"`     // срок хранения в днях
	Compress   bool   `yaml:"compress"`    // сжатие старых логов
}

// DefaultConfig возвращает конфиг по умолчанию
func DefaultConfig() Config {
	return Config{
		Level:      "info",
		Directory:  "./logs",
		Filename:   "app.log",
		MaxSizeMB:  100,
		MaxBackups: 3,
		MaxAgeDays: 30,
		Compress:   true,
	}
}
