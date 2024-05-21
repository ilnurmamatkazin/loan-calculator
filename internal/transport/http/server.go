// Package http - пакет транспортного уровня. Реализует роуты и анализ входных данных в запросе.
package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"loan-calculator/internal/storage"
	"loan-calculator/internal/storage/maps"
	"loan-calculator/pkg/config"
)

// Server - структура http сервера.
type Server struct {
	cfg    *config.Config
	stor   storage.Storager
	Server *http.Server
}

// Init - инициализируем http сервер.
func Init(ctx context.Context) (*Server, error) {
	slog.Info("HTTP сервер создан.")

	cfg, ok := ctx.Value(config.ContextKeyConfig).(*config.Config)
	if !ok {
		return nil, ErrConfig
	}

	s := &Server{
		cfg:  cfg,
		stor: maps.Init(ctx),
		Server: &http.Server{
			Addr:              cfg.HTTP.Host,
			ReadTimeout:       cfg.HTTP.ReadTimeout,
			ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
			WriteTimeout:      cfg.HTTP.WriteTimeout,
		},
	}

	s.Server.Handler = s.NewRouter()

	return s, nil
}

// Run - запускаем http сервер.
func (s *Server) Run() error {
	slog.Info("HTTP сервер запущен.")

	if err := s.Server.ListenAndServe(); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			slog.Info("HTTP сервер остановлен.")
		default:
			return errors.Unwrap(err)
		}
	}

	return nil
}
