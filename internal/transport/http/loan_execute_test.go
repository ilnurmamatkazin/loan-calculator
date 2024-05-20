package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"loan-calculator/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_execute(t *testing.T) {
	s := &Server{
		stor: new(MockStorage),
	}

	//router := s.NewRouter()
	//s.server = httptest.NewServer(),

	type args struct {
		data interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "Positive: тестируем создание записи в кэше",
			args: args{
				data: model.LoanNew{
					ObjectCost:     5000000,
					InitialPayment: 1000000,
					Months:         240,
					Program:        model.Program{Salary: true, Military: false, Base: false},
				},
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name: "Negative: получаем ошибку неверных входных данных",
			args: args{
				data: "",
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Negative: получаем ошибку отсутсвия программ",
			args: args{
				data: model.LoanNew{
					ObjectCost:     5000000,
					InitialPayment: 1000000,
					Months:         240,
					Program:        model.Program{Salary: false, Military: false, Base: false},
				},
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Negative: получаем ошибку множественного выбора программ",
			args: args{
				data: model.LoanNew{
					ObjectCost:     5000000,
					InitialPayment: 1000000,
					Months:         240,
					Program:        model.Program{Salary: true, Military: true, Base: false},
				},
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Negative: получаем ошибку неверного первоночального взноса",
			args: args{
				data: model.LoanNew{
					ObjectCost:     5000000,
					InitialPayment: 100000,
					Months:         240,
					Program:        model.Program{Salary: true, Military: false, Base: false},
				},
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.args.data)

			r := httptest.NewRequest("POST", "/execute", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			err := s.execute(w, r)

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
