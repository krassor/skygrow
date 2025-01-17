package main

import (
	"context"
	"github.com/krassor/skygrow/backend-service-calendar/internal/config"
	"github.com/krassor/skygrow/backend-service-calendar/internal/graceful"
	"github.com/krassor/skygrow/backend-service-calendar/internal/repositories"
	"github.com/krassor/skygrow/backend-service-calendar/internal/services/GoogleService"
	"github.com/krassor/skygrow/backend-service-calendar/internal/services/calendar"
	myHttpServer "github.com/krassor/skygrow/backend-service-calendar/internal/transport/httpServer"
	"github.com/krassor/skygrow/backend-service-calendar/internal/transport/httpServer/handlers"
	"github.com/krassor/skygrow/backend-service-calendar/internal/transport/httpServer/routers"
	"github.com/krassor/skygrow/backend-service-calendar/internal/utils/logger/handlers/slogpretty"
	"log/slog"
	"os"
	"time"
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
		"starting backend-service-calendar",
		slog.String("env", cfg.Env),
		slog.String("version", "0.2"),
	)
	log.Debug("debug messages are enabled")

	googleCalendar := GoogleService.NewGoogleCalendar(log, cfg)
	calendarRepository := repositories.NewCalendarRepository(log, cfg)
	calendarService := calendar.NewCalendarService(log, calendarRepository, googleCalendar)
	calendarHandler := handlers.NewCalendarHandler(log, calendarService)
	router := routers.NewRouter(calendarHandler, cfg.HttpServer.Secret)
	httpServer := myHttpServer.NewHttpServer(log, router, cfg)

	maxSecond := 15 * time.Second
	waitShutdown := graceful.GracefulShutdown(
		context.Background(),
		maxSecond,
		map[string]graceful.Operation{
			"http": func(ctx context.Context) error {
				return httpServer.Shutdown(ctx)
			},
		},
	)

	go httpServer.Listen()
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
