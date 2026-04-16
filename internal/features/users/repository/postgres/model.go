package user_postgres_repository

import "github.com/Kosvu/todoapp-golang/internal/core/domain"

type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func userDomainsFromModels(users []UserModel) []domain.User {
	usersDomain := make([]domain.User, len(users))

	for i, user := range users {
		usersDomain[i] = domain.NewUser(
			user.ID,
			user.Version,
			user.FullName,
			user.PhoneNumber,
		)
	}

	return usersDomain
}
