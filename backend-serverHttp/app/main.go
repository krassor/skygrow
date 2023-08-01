package main

import (
	"context"
	"time"

	"github.com/krassor/skygrow/backend-serverHttp/internal/graceful"
	"github.com/krassor/skygrow/backend-serverHttp/internal/logger"
	"github.com/krassor/skygrow/backend-serverHttp/internal/repositories"
	"github.com/krassor/skygrow/backend-serverHttp/internal/services/bookOrderServices"
	subscriber "github.com/krassor/skygrow/backend-serverHttp/internal/services/subscriberServices"
	telegramBot "github.com/krassor/skygrow/backend-serverHttp/internal/telegram"
	httpServer "github.com/krassor/skygrow/backend-serverHttp/internal/transport/rest-server"
	"github.com/krassor/skygrow/backend-serverHttp/internal/transport/rest-server/handlers"
	"github.com/krassor/skygrow/backend-serverHttp/internal/transport/rest-server/routers"
)

func main() {
	logger.InitLogger()

	repository := repositories.NewRepository()

	
	subscriberService := subscriber.NewSubscriberRepoService(repository)
	bookOrderTgBot := telegramBot.NewBot(subscriberService)

	bookOrderService := bookOrderServices.NewBookOrderService(repository, bookOrderTgBot)

	bookOrderHandler := handlers.NewBookOrderHandler(bookOrderService)
	router := routers.NewRouter(bookOrderHandler)
	httpServer := httpServer.NewHttpServer(router)

	maxSecond := 15 * time.Second
	waitShutdown := graceful.GracefulShutdown(
		context.Background(),
		maxSecond,
		map[string]graceful.Operation{
			"http": func(ctx context.Context) error {
				return httpServer.Shutdown(ctx)
			},
			"tgBotOrders": func(ctx context.Context) error {
				return bookOrderTgBot.Shutdown(ctx)
			},
		},
	)

	go httpServer.Listen()
	go bookOrderTgBot.Update(context.Background(), 60)
	<-waitShutdown
}
