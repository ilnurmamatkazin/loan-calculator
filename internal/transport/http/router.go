package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

// NewRouter - создает таблицу маршрутизации и регестрирует мидлвары.
func (s *Server) NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /execute", handlerFuncWithError(s.execute))
	mux.HandleFunc("GET /cache", handlerFuncWithError(s.cache))

	// Ловим панику
	handler := mwPanicRecovery(mux)

	// Статистика запроса
	handler = mwRequestStatistic(handler)

	return handler
}

// handlerFuncWithError - hendler, возвращающий ошибки.
func handlerFuncWithError(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			var target errorStatus
			status := http.StatusInternalServerError

			if errors.As(err, &target) {
				status = target.status
			}

			var e = struct {
				Error string `json:"error"`
			}{
				Error: err.Error(),
			}

			writeResponse(w, status, e)
		}
	}
}

// Мидлваре по обработке паники в хендлерах.
func mwPanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				slog.Error(string(debug.Stack()))
			}
		}()

		next.ServeHTTP(w, req)
	})
}

// Мидлваре расчета статистики.
func mwRequestStatistic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		sr := &StatusRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}

		next.ServeHTTP(sr, req)

		duration := time.Since(start)
		text := fmt.Sprintf("%s status_code: %d, duration: %d ns\n", start.Format("2006/01/02 15:04:05"), sr.Status, duration)

		if _, err := os.Stdout.WriteString(text); err != nil {
			slog.Error(err.Error())
		}
	})
}

// Функция по формированию ответа клиенту.
func writeResponse(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
