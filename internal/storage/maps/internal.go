package maps

import (
	"context"
	"log/slog"
	"math"
	"reflect"
	"strings"
	"time"

	"loan-calculator/internal/model"
	"loan-calculator/pkg/config"
)

// Так как изначально, ничего не сказано про размер хранилища,
// то для уменьшения аллокаций памяти, инициализируем мапу определенной емкости.
func (st *Storage) new(ctx context.Context) {
	countMapItems := 1024

	// Получаем размер хранилища из конфигурационного файла
	if cfg, ok := ctx.Value(config.ContextKeyConfig).(*config.Config); ok {
		countMapItems = cfg.App.CountMapItems
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	// Инициализируем мапу кредитных программ
	st.mapPrograms = map[string]int{"salary": 8, "military": 9, "base": 10}

	// Инициализируем мапу хранилица
	st.m = make(map[int]model.Loan, countMapItems)

	slog.Info("Мапа проинициализированна.")
}

func (st *Storage) calculationLoan(data model.LoanNew) model.Loan {
	// Так как структура прошла валидацию, то ошибку отсутсвия ставки не проверяем
	rate := st.getRate(data.Program)

	// Ежемесечная процентная ставка
	g := float64(rate) / float64(1200)

	// Сумма ипотечной задолженности
	s := data.ObjectCost - data.InitialPayment

	// Размер ежемесячного платежа
	pm := math.Ceil(float64(s) * g * math.Pow((1+g), float64(data.Months)) / (math.Pow((1+g), float64(data.Months)) - 1))

	return model.Loan{
		Aggregates: model.Aggregates{
			LastPaymentDate: time.Now().AddDate(0, data.Months, 0).Format("2006-01-02"),
			Rate:            rate,
			LoanSum:         s,
			MonthlyPayment:  int(pm),
			Overpayment:     int(pm)*data.Months - s,
		},
		Params: model.Params{
			ObjectCost:     data.ObjectCost,
			InitialPayment: data.InitialPayment,
			Months:         data.Months,
		},
		Program: data.Program,
	}
}

func (st *Storage) getRate(program model.Program) int {
	var (
		fieldName string
		field     reflect.StructField
	)
	v := reflect.ValueOf(program)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field = t.Field(i)
		fieldName = field.Name

		if v.Field(i).Bool() {
			break
		}
	}

	return st.mapPrograms[strings.ToLower(fieldName)]
}
