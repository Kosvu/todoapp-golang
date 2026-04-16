package core_http_middleware

import (
	"context"
	"net/http"
	"time"

	core_logger "github.com/Kosvu/todoapp-golang/internal/core/logger"
	core_http_response "github.com/Kosvu/todoapp-golang/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// просто обьявление константы
const (
	requestIDHeader = "X-Request-ID"
)

// middleware которая получает requestID, а если его нет, то создает
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

/*
middleware которая получает requestID встраивает
его и путь в логер и кладет этот логер в контекст запроса и
просто прокидывает дальше
*/
func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			log = log.With(
				zap.String("requestID", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := context.WithValue(r.Context(), "log", log)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

/*
middleware которая ловит панику
и обрабатывает ее
*/
func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPRespsonseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(
						p,
						"during handle HTTP request got unexpected panic",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

/*
middleware которое логирует вход и сколько времени прошло
*/
func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">>> incoming HTTP request",
				zap.String("http_method", r.Method),
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP request",
				zap.Int("status code", rw.GetStatusCodeOrPanic()),
				zap.Duration("latency", time.Since(before)),
			)
		})
	}
}
