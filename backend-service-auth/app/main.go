package main

import (
	"context"
	"github.com/krassor/skygrow/backend-service-auth/utils/logger/handlers/slogpretty"
	"log/slog"
	"os"
	"time"

	"github.com/krassor/skygrow/backend-service-auth/internal/config"
	"github.com/krassor/skygrow/backend-service-auth/internal/graceful"
	"github.com/krassor/skygrow/backend-service-auth/internal/repositories"
	"github.com/krassor/skygrow/backend-service-auth/internal/services/userServices"
	httpServer "github.com/krassor/skygrow/backend-service-auth/internal/transport/rest-server"
	"github.com/krassor/skygrow/backend-service-auth/internal/transport/rest-server/handlers"
	"github.com/krassor/skygrow/backend-service-auth/internal/transport/rest-server/routers"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting backend-service-auth",
		slog.String("env", cfg.Env),
		slog.String("version", "0.2"),
	)
	log.Debug("debug messages are enabled")

	repository := repositories.NewRepository(cfg)

	userService := userServices.NewUser(repository, cfg.HttpServer.Secret)

	userHandler := handlers.NewUserHandler(userService)

	router := routers.NewRouter(userHandler)
	newHttpServer := httpServer.NewHttpServer(router)

	maxSecond := 15 * time.Second
	waitShutdown := graceful.GracefulShutdown(
		context.Background(),
		maxSecond,
		map[string]graceful.Operation{
			"http": func(ctx context.Context) error {
				return newHttpServer.Shutdown(ctx)
			},
		},
	)

	go newHttpServer.Listen(cfg)
	<-waitShutdown
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
