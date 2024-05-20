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
	server *http.Server
}

// Init - инициализируем http сервер.
func Init(ctx context.Context, stor *maps.Storage) *Server {
	slog.Info("HTTP сервер создан.")

	return &Server{
		cfg:  ctx.Value(config.ContextKeyConfig).(*config.Config),
		stor: stor,
	}
}

// Run - запускаем http сервер.
func (s *Server) Run() error {
	s.server = &http.Server{
		Addr:              s.cfg.HTTP.Host,
		Handler:           s.NewRouter(),
		ReadTimeout:       s.cfg.HTTP.ReadTimeout,
		ReadHeaderTimeout: s.cfg.HTTP.ReadHeaderTimeout,
		WriteTimeout:      s.cfg.HTTP.WriteTimeout,
	}

	slog.Info("HTTP сервер запущен.", "host", s.cfg.HTTP.Host)

	if err := s.server.ListenAndServe(); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			slog.Info("HTTP сервер остановлен.")
		default:
			return errors.Unwrap(err)
		}
	}

	return nil
}

// CloseGracefully - останавливаем http сервер.
func (s *Server) CloseGracefully(ctx context.Context) {
	if err := s.server.Shutdown(ctx); err != nil {
		slog.Error("HTTP shutdown", "error", err)

		if err = s.server.Close(); err != nil {
			slog.Error("HTTP close", "error", err)
		}
	}
}
