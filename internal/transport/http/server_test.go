// Package http - пакет транспортного уровня. Реализует роуты и анализ входных данных в запросе.
package http

import (
	"context"
	"errors"
	"loan-calculator/pkg/config"
	"net/http"
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
	//stor := maps.Init(ctx)

	server, err := Init(ctx)

	// Проверяем результаты
	assert.Nil(t, err)
	assert.NotNil(t, server)
	assert.Equal(t, cfg, server.cfg)
	//assert.Equal(t, stor, server.stor)

	go func() {
		time.Sleep(1 * time.Second)
		server.Server.Shutdown(ctx)
	}()

	err = server.Run()

	switch {
	case err == nil:
		assert.Nil(t, err)
	case errors.Is(err, http.ErrServerClosed):
		assert.Equal(t, err, http.ErrServerClosed)
	default:
		assert.Error(t, err)
	}
}
