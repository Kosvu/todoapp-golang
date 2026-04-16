package core_postgres_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

//Интерфейс который поможет при тестировании

type Pool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Close()

	//время на удержание подключения
	OpTimeout() time.Duration
}

//Cтруктур для пула подлючений

type ConnectionPool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func NewConnectionPool(
	ctx context.Context,
	config Config,
) (*ConnectionPool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Database,
	)

	pgxconfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("create pgxconfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pool ping: %w", err)
	}

	return &ConnectionPool{
		Pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.opTimeout
}
