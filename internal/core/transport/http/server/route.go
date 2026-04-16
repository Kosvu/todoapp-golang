package core_http_server

import "net/http"

/*
Обьявляем структуру нашего машрута, чтобы потом было
легче использовать в mux

Маршрут состоит из метода, пути и самого хендера
*/

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

/*
Конструктор маршрута
*/
func NewRoute(
	method string,
	path string,
	handler http.HandlerFunc,
) Route {
	return Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}
