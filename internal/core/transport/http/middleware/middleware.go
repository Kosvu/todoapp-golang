package core_http_middleware

import "net/http"

/*
Создаем алиас middleware чтобы код смотрелся красиво
То есть когда мы так пишем, то Middleware и
func(http.Handler) http.Handler взаимозаменяемы
*/

type Middleware func(http.Handler) http.Handler

// Создаем цепочку middleware
// Идем с конца потому что нам очень важна последовательность
func ChainMiddleware(
	h http.Handler,
	m ...Middleware,
) http.Handler {
	if len(m) == 0 {
		return h
	}

	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}
