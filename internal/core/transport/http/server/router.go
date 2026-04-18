package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/Kosvu/todoapp-golang/internal/core/transport/http/middleware"
)

/*
Тут мы делаем версионирование нашему API
*/

/*
Создаем тип на основе стринг, чтобы не путаться
*/
type ApiVersion string

var (
	ApiVersion1 = ApiVersion("v1")
	ApiVersion2 = ApiVersion("v2")
	ApiVersion3 = ApiVersion("v3")
)

/*
Структура маршрутизатор с версиями
Встраиваем ServeMux, чтобы использовать регистрацию через Handle
и Еще пишем после нашей версии
*/
type APIVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
	middleware []core_http_middleware.Middleware
}

/*
Конструктор
*/
func NewAPIVersionRouter(
	apiVersion ApiVersion,
	middleware ...core_http_middleware.Middleware,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
		middleware: middleware,
	}
}

/*
Регистрация роутов, которая принимает обьект роут
склеивает метод, путь и потом еще Handler
*/
func (r *APIVersionRouter) RegisterRoute(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.WithMiddleware())
	}
}

func (r *APIVersionRouter) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r,
		r.middleware...,
	)
}
