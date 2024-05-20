// Package maps - пакет, реализующий хранение даннных в RAM.
package maps

import (
	"context"
	"log/slog"
	"slices"
	"sync"

	"loan-calculator/internal/model"
)

// Storage - объект который реализует интерфейс хранилища.
type Storage struct {
	m           map[int]model.Loan
	mapPrograms map[string]int
	mutex       sync.RWMutex
	autoKey     int // счетчик автоинкремента
}

var (
	once sync.Once
	stor *Storage
)

// Init - функция инициализации хранилища.
func Init(ctx context.Context) *Storage {
	once.Do(func() {
		stor = &Storage{}
		slog.Info("Создан Storage.")

		// Создаем хранилице в RAM
		stor.new(ctx)
	})

	return stor
}

// Execute - функция добавления нового значения в мапу.
func (st *Storage) Execute(data model.LoanNew) model.Loan {
	loan := st.calculationLoan(data)

	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.autoKey++
	loan.ID = st.autoKey
	st.m[st.autoKey] = loan

	return loan
}

// Cache - функция возвращает кэш в виде слайса.
func (st *Storage) Cache() []model.Loan {
	st.mutex.RLock()
	defer st.mutex.RUnlock()

	// Сортируем ключи мапы
	keys := make([]int, 0, len(st.m))
	for k := range st.m {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	// Выводим финальный слайс в отсортированном виде
	loans := make([]model.Loan, 0, len(st.m))
	for _, k := range keys {
		loans = append(loans, st.m[k])
	}

	return loans
}
