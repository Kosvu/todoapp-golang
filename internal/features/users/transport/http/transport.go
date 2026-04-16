package users_transport_http

import (
	"context"
	"net/http"

	"github.com/Kosvu/todoapp-golang/internal/core/domain"
	core_http_server "github.com/Kosvu/todoapp-golang/internal/core/transport/http/server"
)

// Структура транспорта UserHTTPHandler и показываем что она зависит от сервиса
type UsersHTTPHandler struct {
	usersService UsersService
}

// интерфейс этого самого сервиса, чтобы мы могли скрыть работу самого сервиса
type UsersService interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		id int,
	) error

	PatchUser(
		ctx context.Context,
		id int,
		patch domain.UserPatch,
	) (domain.User, error)
}

// конструктор создания transport слоя
func NewUsersHTTPHandler(
	usersService UsersService,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

// метод который возвращает все ендпоинты на этом транспортном слое
func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: h.PatchUser,
		},
	}
}
