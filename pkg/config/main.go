// Package config - пакет формирования конфигурации.
package config

import (
	"errors"
	"log/slog"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	configPath = "./cmd/config.yml"
	once       sync.Once
	cfg        *Config
)

// ContextKeyConfig - переменная необходима для чтения/записи в контекст.
var ContextKeyConfig contextKeyConfig

// New - инициализируем конфигурацию, загружаем данные из конфигурационного файла.
func New() (*Config, error) {
	var (
		err      error
		yamlFile []byte
	)

	once.Do(func() {
		if yamlFile, err = os.ReadFile(configPath); err != nil {
			return
		}

		cfg = &Config{}
		if err = yaml.Unmarshal(yamlFile, cfg); err != nil {
			return
		}

		slog.Info("Конфигурация загружена.")
	})

	return cfg, errors.Unwrap(err)
}
