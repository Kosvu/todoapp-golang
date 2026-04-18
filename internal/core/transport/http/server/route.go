package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/Kosvu/todoapp-golang/internal/core/transport/http/middleware"
)

/*
Обьявляем структуру нашего машрута, чтобы потом было
легче использовать в mux

Маршрут состоит из метода, пути и самого хендера
*/

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []core_http_middleware.Middleware
}

func (r *Route) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r.Handler,
		r.Middleware...,
	)
}

/*
Конструктор маршрута
*/

//Пока не нужен
// func NewRoute(
// 	method string,
// 	path string,
// 	handler http.HandlerFunc,
// ) Route {
// 	return Route{
// 		Method:  method,
// 		Path:    path,
// 		Handler: handler,
// 	}
// }
