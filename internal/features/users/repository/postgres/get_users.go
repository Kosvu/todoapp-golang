package user_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Kosvu/todoapp-golang/internal/core/domain"
)

func (r *UserRepository) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, full_name, phone_number
	FROM todoapp.users
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2 
	`

	rows, err := r.pool.Query(
		ctx,
		query,
		limit,
		offset,
	)

	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()

	var usersModels []UserModel

	for rows.Next() {
		var userModel UserModel

		err := rows.Scan(
			&userModel.ID,
			&userModel.Version,
			&userModel.FullName,
			&userModel.PhoneNumber,
		)

		if err != nil {
			return nil, fmt.Errorf("scan users: %w", err)
		}

		usersModels = append(usersModels, userModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows next: %w", err)
	}

	userDomains := userDomainsFromModels(usersModels)

	return userDomains, nil
}
