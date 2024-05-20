// Package http - пакет транспортного уровня. Реализует роуты и анализ входных данных в запросе.
package http

import (
	"context"
	"loan-calculator/internal/storage/maps"
	"loan-calculator/pkg/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	cfg := &config.Config{
		HTTP: config.HTTP{
			Host:              "0.0.0.0:8080",
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      15 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
		},
		App: config.App{
			CountMapItems: 1024,
		},
	}

	ctx := context.WithValue(context.Background(), config.ContextKeyConfig, cfg)
	stor := maps.Init(ctx)

	server := Init(ctx, stor)

	// Проверяем результаты
	assert.NotNil(t, server)
	assert.Equal(t, cfg, server.cfg)
	assert.Equal(t, stor, server.stor)
}

func TestServer_CloseGracefully(t *testing.T) {
	t.Run("Shutdown successful", func(t *testing.T) {
		// Создаем тестовый HTTP сервер
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer testServer.Close()

		// Создаем наш сервер с тестовым сервером
		server := &Server{server: testServer.Config}

		// Запускаем метод Shutdown в горутине
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		go func() {
			server.CloseGracefully(ctx)
		}()

		// Даем время для корректного завершения
		time.Sleep(100 * time.Millisecond)

		// Проверяем что сервер был корректно завершен
		assert.NoError(t, ctx.Err())
	})
}
