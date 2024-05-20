// Package maps - пакет, реализующий хранение даннных в RAM.
package maps

import (
	"context"
	"loan-calculator/internal/model"
	"reflect"
	"testing"
	"time"
)

func TestStorage_Cache(t *testing.T) {
	type fields struct {
		m map[int]model.Loan
	}
	tests := []struct {
		name   string
		fields fields
		want   []model.Loan
	}{
		{
			name: "Тестируем функцию cache",
			fields: fields{
				m: make(map[int]model.Loan),
			},
			want: make([]model.Loan, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &Storage{
				m: tt.fields.m,
			}

			if got := st.Cache(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.Cache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Execute(t *testing.T) {
	type args struct {
		data model.LoanNew
	}
	tests := []struct {
		name string
		args args
		want model.Loan
	}{
		{
			name: "Тестируем функцию Execute",
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
				ID: 1,
			},
		},
	}

	st := &Storage{
		m:           make(map[int]model.Loan),
		mapPrograms: map[string]int{"salary": 8, "military": 9, "base": 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := st.Execute(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInit(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *Storage
	}{
		{
			name: "Тестируем функцию создания хранилища",
			args: args{
				ctx: context.Background(),
			},
			want: &Storage{
				m:           make(map[int]model.Loan),
				mapPrograms: map[string]int{"salary": 8, "military": 9, "base": 10},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Init(tt.args.ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}

			if reflect.DeepEqual(got.m, nil) {
				t.Errorf("Init() storage not init")
			}

			if reflect.DeepEqual(got.mapPrograms, nil) {
				t.Errorf("Init() programs not init")
			}
		})
	}
}
