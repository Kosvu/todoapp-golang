package core_postgres_pool

import "errors"

var (
	ErrNoRows = errors.New("no rows")

	// нарушение связи с внешним ключом, если например его попросту нет
	ErrViolatesForeignKey = errors.New("violates foreign key")

	// какая-то неизвестная ошибка
	ErrUnknown = errors.New("unknown")
)
