package main

import (
	"app/main.go/internal/config"
	"app/main.go/internal/graceful"
	"app/main.go/internal/utils/logger/handlers/slogpretty"
	"context"
	"log/slog"
	"os"
	"time"
	telegramBot "app/main.go/internal/telegram"
	deepseek "app/main.go/internal/deepseek"
	inMemory "app/main.go/internal/cache/inMemory"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var Version = "dev"

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting backend-miniaps-bot",
		slog.String("env", cfg.Env),
		slog.String("version", Version),
	)
	log.Debug("debug messages are enabled")

	inMemoryCache := inMemory.NewInMemoryRepository()
	AIBot := deepseek.NewClient(log, cfg, inMemoryCache)
	tgBot := telegramBot.New(log, cfg, AIBot)

	maxSecond := 15 * time.Second
	waitShutdown := graceful.GracefulShutdown(
		context.Background(),
		maxSecond,
		map[string]graceful.Operation{
			// "http": func(ctx context.Context) error {
			// 	return httpServer.Shutdown(ctx)
			// },
			"tgBot": func(ctx context.Context) error {
				return tgBot.Shutdown(ctx)
			},
		},
		log,
	)
	go tgBot.Update(60)
	// go httpServer.Listen()
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
