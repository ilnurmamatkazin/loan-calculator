// Package main - главный файл приложения
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"loan-calculator/internal/transport/http"
	"loan-calculator/pkg/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Настраиваем логер, что бы данные отдавались в формате json
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Помещаем объект конфигурации в контекст
	ctx = initConfig(ctx)

	// Инициализируем http сервер
	httpServer, err := http.Init(ctx)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// Ловим сигнал завершения работы сервиса
	chSignal := make(chan os.Signal, 2)
	signal.Notify(chSignal, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	tmRun := time.NewTimer(0)
	isContinue := true

	var wg sync.WaitGroup

	for isContinue {
		select {
		case <-chSignal:
			slog.Info("Начало остановки сервиса...")

			// Шлём во все горутины сигнал завершения работы
			cancel()

			// Завершаем работу http сервера
			if err := httpServer.Server.Shutdown(ctx); err != nil {
				slog.Error("HTTP shutdown", "error", err)
			}

			isContinue = false

		case <-tmRun.C:
			wg.Add(1)
			go func() {
				defer func() {
					if rec := recover(); rec != nil {
						slog.Error("%v / %v", rec, string(debug.Stack()))
					}

					wg.Done()

					// Если по какой-то причине завершается горутина, шлем сигнал завершения приложения
					chSignal <- syscall.SIGINT
				}()

				var err error

				// Запускаем http сервер
				if err = httpServer.Run(); err != nil {
					slog.Error(err.Error())
					return
				}
			}()
		}
	}

	wg.Wait()
	slog.Info("Сервис выключен.")
}

func initConfig(ctx context.Context) context.Context {
	// Инициализируем конфигурацию
	conf, err := config.New()
	if err != nil {
		slog.Error("Чтение конфигурации", "error", err)
		os.Exit(1)
	}

	// Записываем конфигурацию в контекст
	ctx = context.WithValue(ctx, config.ContextKeyConfig, conf)

	return ctx
}
