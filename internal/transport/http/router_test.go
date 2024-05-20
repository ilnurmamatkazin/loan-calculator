package http

import (
	"loan-calculator/internal/model"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Заглушка для хранилища
type MockStorage struct {
	mock.Mock
	m map[int]struct{}
}

func (m *MockStorage) Execute(body model.LoanNew) model.Loan {
	return model.Loan{}
}

func (m *MockStorage) Cache() []model.Loan {
	if m.m == nil {
		return nil
	}

	return make([]model.Loan, 1)
}

func TestServer_NewRouter(t *testing.T) {
	server := &Server{
		stor: &MockStorage{m: make(map[int]struct{})},
	}

	router := server.NewRouter()

	tests := []struct {
		name     string
		method   string
		path     string
		statuses []int
		wantErr  bool
	}{
		{
			name:     "Positive: тестируем возврат нужного статуса из заявленного перечня для хендлера /execute",
			method:   "POST",
			path:     "/execute",
			statuses: []int{http.StatusOK, http.StatusBadRequest, http.StatusInternalServerError},
			wantErr:  false,
		},
		{
			name:     "Negative: тестируем возврат нужного статуса из заявленного перечня для хендлера /execute",
			method:   "GET",
			path:     "/execute",
			statuses: []int{http.StatusOK, http.StatusBadRequest, http.StatusInternalServerError},
			wantErr:  true,
		},
		{
			name:     "Positive: тестируем возврат нужного статуса из заявленного перечня для хендлера /cache",
			method:   "GET",
			path:     "/cache",
			statuses: []int{http.StatusOK, http.StatusBadRequest, http.StatusInternalServerError},
			wantErr:  false,
		},
		{
			name:     "Negative: тестируем возврат нужного статуса из заявленного перечня для хендлера /cache",
			method:   "POST",
			path:     "/cache",
			statuses: []int{http.StatusOK, http.StatusBadRequest, http.StatusInternalServerError},
			wantErr:  true,
		},
		{
			name:     "Positive: тестируем возврат статуса 404",
			method:   "GET",
			path:     "/404",
			statuses: []int{http.StatusNotFound},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.True(t, tt.wantErr != slices.Contains(tt.statuses, rr.Code))
		})
	}
}
