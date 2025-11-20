package main

import (
	"app/main.go/internal/config"
	"app/main.go/internal/graceful"
	"app/main.go/internal/mail"
	"app/main.go/internal/openrouter"
	"app/main.go/internal/pdf"
	"app/main.go/internal/repositories"
	telegramBot "app/main.go/internal/telegram"
	"app/main.go/internal/transport/httpServer"
	"app/main.go/internal/transport/httpServer/handlers"
	"app/main.go/internal/transport/httpServer/routers"
	"app/main.go/internal/utils/logger/handlers/slogpretty"
	"context"
	"log/slog"
	"os"
	"time"
	//openai "app/main.go/internal/openai"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var Version = "0.1"

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting backend-service-profreport",
		slog.String("env", cfg.Env),
		slog.String("version", Version),
	)

	err := cfg.ReadPromptFromFile()
	if err != nil {
		log.Error(
			"main read prompt error",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	tgBot := telegramBot.New(log, cfg)

	// инициализируем репозиторий для работы с БД
	repository := repositories.NewCalendarRepository(log, cfg)

	mailService := mail.NewMailer(log, cfg)
	pdfService := pdf.New(log, cfg, mailService)
	openRouterService := openrouter.NewClient(log, cfg, pdfService)
	questionnaireHandler := handlers.NewQuestionnaireHandler(
		log,
		cfg,
		openRouterService,
		repository,
	)

	router := routers.NewRouter(questionnaireHandler, "TODO")
	httpServer := httpServer.NewHttpServer(log, router, cfg)

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
			"httpServer": func(ctx context.Context) error {
				return httpServer.Shutdown(ctx)
			},
			"Mailer": func(ctx context.Context) error {
				return mailService.Shutdown(ctx)
			},
			"Pdf service": func(ctx context.Context) error {
				return pdfService.Shutdown(ctx)
			},
			"Openrouter service": func(ctx context.Context) error {
				return openRouterService.Shutdown(ctx)
			},
		},
		log,
	)
	go tgBot.Update(60)
	go mailService.Start()
	go pdfService.Start()
	go openRouterService.Start()
	go httpServer.Listen()
	// err := mailService.AddJob("krassor86@yandex.ru", "test proffreport", "body2")
	// if err != nil {
	// 	log.Error(
	// 		"main send mail error",
	// 		slog.String("error", err.Error()),
	// 	)
	// }
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
		// log = slog.New(
		// 	slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		// )
		log = setupPrettySlogProd()
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

func setupPrettySlogProd() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
