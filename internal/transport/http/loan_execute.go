package http

import (
	"encoding/json"
	"loan-calculator/internal/model"
	"math"
	"net/http"
	"reflect"
)

func (s *Server) execute(w http.ResponseWriter, r *http.Request) (err error) {
	var body model.LoanNew
	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return errorStatus{error: err, status: http.StatusBadRequest}
	}
	defer func() {
		if errDefer := r.Body.Close(); errDefer != nil {
			err = errorStatus{error: errDefer, status: http.StatusInternalServerError}
			return
		}
	}()

	if status, err := validateExecute(body); err != nil {
		return errorStatus{error: err, status: status}
	}

	writeResponse(w, http.StatusOK, s.stor.Execute(body))

	return nil
}

func validateExecute(body model.LoanNew) (int, error) {
	// Количество переданных клиентом программ
	fieldCount := 0

	// Получаем значение структуры через reflect.ValueOf
	v := reflect.ValueOf(body.Program)

	// Определяем количество выбранных пользователем программ
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			fieldValue := v.Field(i)

			if fieldValue.Kind() == reflect.Bool {
				if fieldValue.Bool() {
					fieldCount++
				}
			} else {
				return http.StatusInternalServerError, ErrIsNotTypeBool
			}
		}
	} else {
		return http.StatusInternalServerError, ErrIsNotStruct
	}

	// Проверяем значение количесва программ, выбранных пользователем
	switch fieldCount {
	case 0:
		return http.StatusBadRequest, ErrChooseProgram

	case 1:
		if body.InitialPayment < int(math.Ceil(float64(body.ObjectCost)*0.2)) {
			return http.StatusBadRequest, ErrInitialPayment
		}

	default:
		return http.StatusBadRequest, ErrChooseOnlyProgram
	}

	return http.StatusOK, nil
}
