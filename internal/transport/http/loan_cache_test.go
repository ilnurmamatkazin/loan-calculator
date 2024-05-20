package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_cache(t *testing.T) {
	//router := s.NewRouter()
	//s.server = httptest.NewServer(),

	tests := []struct {
		name       string
		ms         *MockStorage
		wantErr    bool
		wantStatus int
	}{
		{
			name:       "Positive: тестируем выдачу кэша",
			ms:         &MockStorage{m: make(map[int]struct{})},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Negative: получаем ошибку пустого кеша",
			ms:         &MockStorage{},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				stor: tt.ms,
			}

			r := httptest.NewRequest("GET", "/cache", nil)
			w := httptest.NewRecorder()

			err := s.cache(w, r)

			var target errorStatus
			status := http.StatusOK

			if errors.As(err, &target) {
				status = target.status
			}

			assert.EqualValuesf(t, err != nil, tt.wantErr, "%v", err)
			assert.EqualValuesf(t, status, tt.wantStatus, "%v", err)
		})
	}
}
