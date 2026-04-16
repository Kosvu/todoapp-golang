package core_http_server

import (
	"fmt"
	"net/http"
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
}

/*
Конструктор
*/
func NewAPIVersionRouter(
	apiVersion ApiVersion,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

/*
Регистрация роутов, которая принимает обьект роут
склеивает метод, путь и потом еще Handler
*/
func (r *APIVersionRouter) RegisterRoute(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.Handler)
	}
}
