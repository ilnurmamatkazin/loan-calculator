// Package config - пакет формирования конфигурации.
package config

import (
	"errors"
	"flag"
	"log/slog"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	conf = flag.String("config", "./config.yml", "Path to configuration file")
	once sync.Once
	cfg  *Config
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
		flag.Parse()

		if yamlFile, err = os.ReadFile(*conf); err != nil {
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
