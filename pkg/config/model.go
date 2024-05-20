package config

import "time"

type contextKeyConfig struct{}

// Config - структура конфигурации.
type Config struct {
	HTTP HTTP `yaml:"http"`
	App  App  `yaml:"app"`
}

// HTTP - секция настройки http сервера.
type HTTP struct {
	Host              string        `yaml:"host"`
	ReadTimeout       time.Duration `yaml:"read-timeout"`
	WriteTimeout      time.Duration `yaml:"write-timeout"`
	ReadHeaderTimeout time.Duration `yaml:"read-header-timeout"`
}

// App - секция настройки приложения.
type App struct {
	CountMapItems int `yaml:"count-map-items"`
}
