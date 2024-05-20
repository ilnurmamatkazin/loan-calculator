// Package storage - описывает интерфейс хранилища.
package storage

import (
	"loan-calculator/internal/model"
)

// Storager - интерфейс хранилища.
type Storager interface {
	Execute(data model.LoanNew) model.Loan
	Cache() []model.Loan
}
