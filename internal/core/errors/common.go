package core_errors

import "errors"

// обьявление ошибок, которые у нас встречаются в коде в дальнейшем

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrConflict        = errors.New("conflict")
)
