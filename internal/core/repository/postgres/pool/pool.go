package core_postgres_pool

import (
	"context"
	"time"
)

//Интерфейс который поможет при тестировании

type Pool interface {
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, sql string, arguments ...any) (CommandTag, error)
	Close()

	//время на удержание подключения
	OpTimeout() time.Duration
}

// Кастомный типы, чтобы не привязываться к pgx
// Просто смотрим какие методы нам требуются
// И после этого вписываем в наши интерфейсы
type Rows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...any) error
}

type Row interface {
	Scan(dest ...any) error
}

type CommandTag interface {
	RowsAffected() int64
}
