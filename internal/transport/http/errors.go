package http

import "errors"

type errorStatus struct {
	error
	status int
}

var (
	// ErrChooseProgram - не выбрали ни одну программу.
	ErrChooseProgram = errors.New("choose program")

	// ErrInitialPayment - ошибочное значение первоначального взноса.
	ErrInitialPayment = errors.New("the initial payment should be more")

	// ErrChooseOnlyProgram - выбранно более одной программы.
	ErrChooseOnlyProgram = errors.New("choose only 1 program")

	// ErrIsNotStruct - ошибка входных данных.
	ErrIsNotStruct = errors.New("program value is not a struct")

	// ErrIsNotTypeBool - ошибка входных данных.
	ErrIsNotTypeBool = errors.New("field program is not of type bool")

	// ErrEmptyCache - пустой кеш.
	ErrEmptyCache = errors.New("empty cache")
)
