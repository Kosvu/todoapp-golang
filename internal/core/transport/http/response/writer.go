package core_http_response

import "net/http"

/*
Тут находится кастомный writer
который еще и возвращает statuscode
чтобы мы могли его использовать
*/

// переменная если статус кода нет
var (
	StatusCodeUninitialized = -1
)

// Встраиваем обычный writer и еще добавляем поле statuscode
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// конструктор
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     StatusCodeUninitialized,
	}
}

// Метод который записывает статус код и сохраняет его в поле
func (rw *ResponseWriter) WriteHeader(statuscode int) {
	rw.ResponseWriter.WriteHeader(statuscode)
	rw.statusCode = statuscode
}

// метод который вызывает панику если у нас нет статус кода
func (rw *ResponseWriter) GetStatusCodeOrPanic() int {
	if rw.statusCode == StatusCodeUninitialized {
		panic("no status code set")
	}

	return rw.statusCode
}
