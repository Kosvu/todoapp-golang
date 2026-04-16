package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/Kosvu/todoapp-golang/internal/core/logger"
	core_http_middleware "github.com/Kosvu/todoapp-golang/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

/*
тут реализовывается запуск нашего сервера
все работает через стандратный mux

Структрура принимает config, mux, log и список middleware, чтобы обернуть хендлеры
*/
type HTTPServer struct {
	mux    *http.ServeMux
	config Config
	log    *core_logger.Logger

	middleware []core_http_middleware.Middleware
}

/*
Конструктор

Мне стоит обратить внимание что для инициализации mux
автор пишет http.NewServeMux()
*/
func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

/*
Основная регистрация роутов с версионизацией
*/
func (h *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		h.mux.Handle(
			// + "/" означает все что начинается с этого префикса
			prefix+"/",
			http.StripPrefix(prefix, router),
		)
	}
}

/*
Запуск самого сервера
*/
func (h *HTTPServer) Run(ctx context.Context) error {

	// Оборачиваем хендлер
	mux := core_http_middleware.ChainMiddleware(h.mux, h.middleware...)

	// передаем из порт из конфига и хендеры внутри переменной mux
	server := &http.Server{
		Addr:    h.config.Addr,
		Handler: mux,
	}

	// создаем буферизированный канал, чтобы не блокировать запись
	ch := make(chan error, 1)

	// в отдельной горутине запускаем сервер
	// потому что хотим слушать ошибки
	go func() {
		defer close(ch)

		h.log.Warn("start HTTP server", zap.String("addr", h.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	// смотри что пришло в канал первым
	// исходя из этого закрываем сервер или вызываем ошибки
	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and server HTTP: %w", err)
		}
	case <-ctx.Done():
		h.log.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			h.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		h.log.Warn("HTTP server stopped")
	}

	return nil
}
