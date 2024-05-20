package maps

import (
	"context"
	"loan-calculator/internal/model"
	"loan-calculator/pkg/config"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStorage_getRate(t *testing.T) {
	type args struct {
		program model.Program
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Testing salary",
			args: args{
				program: model.Program{Salary: true},
			},
			want: 8,
		},
		{
			name: "Testing military",
			args: args{
				program: model.Program{Military: true},
			},
			want: 9,
		},
		{
			name: "Testing base",
			args: args{
				program: model.Program{Base: true},
			},
			want: 10,
		},
	}

	st := &Storage{
		mapPrograms: map[string]int{"salary": 8, "military": 9, "base": 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := st.getRate(tt.args.program); got != tt.want {
				t.Errorf("Storage.getRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_calculationLoan(t *testing.T) {
	type args struct {
		data model.LoanNew
	}

	tests := []struct {
		name string
		args args
		want model.Loan
	}{
		{
			name: "Testing salary",
			args: args{
				data: model.LoanNew{
					ObjectCost:     5000000,
					InitialPayment: 1000000,
					Months:         240,
					Program:        model.Program{Salary: true, Military: false, Base: false},
				},
			},
			want: model.Loan{
				Aggregates: model.Aggregates{
					LastPaymentDate: time.Now().AddDate(0, 240, 0).Format("2006-01-02"),
					Rate:            8,
					LoanSum:         4000000,
					MonthlyPayment:  33458,
					Overpayment:     4029920,
				},
				Params: model.Params{
					ObjectCost:     5000000,
					InitialPayment: 1000000,
					Months:         240,
				},
				Program: model.Program{
					Salary:   true,
					Military: false,
					Base:     false,
				},
			},
		},
	}

	st := &Storage{
		mapPrograms: map[string]int{"salary": 8, "military": 9, "base": 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := st.calculationLoan(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.calculationLoan() = %v, want %v", got, tt.want)
			}
		})
	}
}

type ContextKey string

const ContextKeyConfig ContextKey = "config"

func TestStorage_new(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Testing new storage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			cfg.App.CountMapItems = 1024

			st := &Storage{}
			ctx := context.WithValue(context.Background(), ContextKeyConfig, cfg)

			st.new(ctx)

			// Проверка результатов
			assert.Equal(t, 3, len(st.mapPrograms))
			assert.Equal(t, 0, len(st.m))
			assert.Equal(t, 8, st.mapPrograms["salary"])
			assert.Equal(t, 9, st.mapPrograms["military"])
			assert.Equal(t, 10, st.mapPrograms["base"])
		})
	}
}
