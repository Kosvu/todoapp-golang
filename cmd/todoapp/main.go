package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/Kosvu/todoapp-golang/internal/core/logger"
	core_pgx_pool "github.com/Kosvu/todoapp-golang/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/Kosvu/todoapp-golang/internal/core/transport/http/middleware"
	core_http_server "github.com/Kosvu/todoapp-golang/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/Kosvu/todoapp-golang/internal/features/tasks/repository/postgres"
	tasks_service "github.com/Kosvu/todoapp-golang/internal/features/tasks/service"
	tasks_transport_http "github.com/Kosvu/todoapp-golang/internal/features/tasks/transport/http"
	user_postgres_repository "github.com/Kosvu/todoapp-golang/internal/features/users/repository/postgres"
	users_service "github.com/Kosvu/todoapp-golang/internal/features/users/service"
	users_transport_http "github.com/Kosvu/todoapp-golang/internal/features/users/transport/http"
	"go.uber.org/zap"
)

var (
	timeZone = time.UTC
)

func main() {

	time.Local = timeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGINT,
	)
	defer cancel()
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to ini application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", timeZone))

	logger.Debug("initializing connection pool")
	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)

	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "user"))

	userRepository := user_postgres_repository.NewUserRepository(pool)
	usersService := users_service.NewUserService(userRepository)
	userTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))

	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTaskService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoute(userTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoute(tasksTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouterV1)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("http server run error", zap.Error(err))
	}

}
