package http

import "net/http"

// StatusRecorder - структура обертка над стандартным ResponseWriter.
// Необходима для перехвата значения статуса.
type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

// WriteHeader - перехватываем статус.
func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}
