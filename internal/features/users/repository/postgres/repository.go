package user_postgres_repository

import core_postgres_pool "github.com/Kosvu/todoapp-golang/internal/core/repository/postgres/pool"

type UserRepository struct {
	pool core_postgres_pool.Pool
}

func NewUserRepository(
	pool core_postgres_pool.Pool,
) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}
