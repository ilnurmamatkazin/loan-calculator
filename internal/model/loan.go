// Package model - файл содержит струкрутры, описывающие кредит.
// Предполагается, что сервис будет запущен на 64 разрядной системе
// и тип int для полей object_cost и initial_payment будет иметь соответсвующий размер
package model

// Loan - структура кредита.
type Loan struct {
	Aggregates Aggregates `json:"aggregates"`
	Params     Params     `json:"params"`
	Program    Program    `json:"program"`
	ID         int        `json:"id"`
}

// LoanNew - структура для создания нового кредита.
type LoanNew struct {
	ObjectCost     int     `json:"object_cost"`
	InitialPayment int     `json:"initial_payment"`
	Months         int     `json:"months"`
	Program        Program `json:"program"`
}

// Params - параметры кредита.
type Params struct {
	ObjectCost     int `json:"object_cost"`
	InitialPayment int `json:"initial_payment"`
	Months         int `json:"months"`
}

// Program - структура кредитной программы.
type Program struct {
	Salary   bool `json:"salary,omitempty"`
	Military bool `json:"military,omitempty"`
	Base     bool `json:"base,omitempty"`
}

// Aggregates - расчетные параметры кредита.
type Aggregates struct {
	LastPaymentDate string `json:"last_payment_date"`
	Rate            int    `json:"rate"`
	LoanSum         int    `json:"loan_sum"`
	MonthlyPayment  int    `json:"monthly_payment"`
	Overpayment     int    `json:"overpayment"`
}

// {
// 	"id": 0, // id расчета в кэше
// 	"params": {
// 	   "object_cost": 5000000,
// 	   "initial_payment": 1000000,
// 	   "months": 240
// 	},
// 	"program": {
// 	   "salary": true
// 	},
// 	"aggregates": {
// 	   "rate": 8,
// 	   "loan_sum": 4000000,
// 	   "monthly_payment": 33458,
// 	   "overpayment": 4029920,
// 	   "last_payment_date": "2044-02-18"
// 	}
//  },
