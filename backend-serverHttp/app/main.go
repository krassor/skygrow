package main

import (
	"context"
	"time"

	"github.com/krassor/skygrow/backend-serverHttp/internal/graceful"
	"github.com/krassor/skygrow/backend-serverHttp/internal/logger"
	"github.com/krassor/skygrow/backend-serverHttp/internal/repositories"
	services "github.com/krassor/skygrow/backend-serverHttp/internal/services/devices"
	fetcher "github.com/krassor/skygrow/backend-serverHttp/internal/services/fetcher"
	subscriber "github.com/krassor/skygrow/backend-serverHttp/internal/services/subscriberServices"
	telegramBot "github.com/krassor/skygrow/backend-serverHttp/internal/telegram"
	httpServer "github.com/krassor/skygrow/backend-serverHttp/internal/transport/rest-sever"
	"github.com/krassor/skygrow/backend-serverHttp/internal/transport/rest-sever/handlers"
	"github.com/krassor/skygrow/backend-serverHttp/internal/transport/rest-sever/routers"
)

func main() {
	logger.InitLogger()

	repository := repositories.NewRepository()

	deviceRepoService := services.NewDeviceRepoService(repository)
	subscriberService := subscriber.NewSubscriberRepoService(repository)

	deviceHandler := handlers.NewDeviceHandler(deviceRepoService)
	deviceRouter := routers.NewDeviceRouter(deviceHandler)
	deviceHttpServer := httpServer.NewHttpServer(deviceRouter)

	deviceTgBot := telegramBot.NewBot(deviceRepoService, subscriberService)

	fetcherDevice := fetcher.NewDeviceFetcher(deviceRepoService, deviceTgBot)

	maxSecond := 15 * time.Second
	waitShutdown := graceful.GracefulShutdown(
		context.Background(),
		maxSecond,
		map[string]graceful.Operation{
			"http": func(ctx context.Context) error {
				return deviceHttpServer.Shutdown(ctx)
			},
			// "tgBot": func(ctx context.Context) error {
			// 	return deviceTgBot.Shutdown(ctx)
			// },
		},
	)

	go deviceHttpServer.Listen()
	go deviceTgBot.Update(context.Background(), 60)
	go fetcherDevice.Start(context.Background())
	<-waitShutdown
}
