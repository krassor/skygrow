package main

import (
	// "context"
	// "fmt"
	// "time"

	// "github.com/krassor/skygrow/backend-miniaps-bot/internal/config"
	// "github.com/krassor/skygrow/backend-miniaps-bot/internal/graceful"
	// "github.com/krassor/skygrow/backend-miniaps-bot/internal/logger"
	// "github.com/krassor/skygrow/backend-miniaps-bot/internal/openai"
	// "github.com/krassor/skygrow/backend-miniaps-bot/internal/repository"
	// telegramBot "github.com/krassor/skygrow/backend-miniaps-bot/internal/telegram"
	// "github.com/krassor/skygrow/backend-miniaps-bot/internal/transport/broker/redisBroker"
	// "github.com/krassor/skygrow/backend-miniaps-bot/internal/transport/httpServer"
	// "github.com/krassor/skygrow/backend-miniaps-bot/internal/transport/httpServer/routers"
)

var Version = "dev"

func main() {
	//1233
	// fmt.Printf("Version: %s\n", Version)

	// logger.InitLogger()

	// appConfig := config.InitConfig()
	// repo := repository.NewMessageRepository()

	// broker := redisBroker.NewRedisClient()

	// gptBot := openai.NewGPTBot(appConfig, repo, broker)
	// tgBot := telegramBot.NewBot(appConfig, broker)

	// botRouter := routers.NewBotRouter()
	// newHttpServer := httpServer.NewHttpServer(botRouter)

	// maxSecond := 15 * time.Second
	// waitShutdown := graceful.GracefulShutdown(
	// 	context.Background(),
	// 	maxSecond,
	// 	map[string]graceful.Operation{
	// 		"tgBot": func(ctx context.Context) error {
	// 			return tgBot.Shutdown(ctx)
	// 		},
	// 		"newHttpServer": func(ctx context.Context) error {
	// 			return newHttpServer.Shutdown(ctx)
	// 		},
	// 		"gptBot": func(ctx context.Context) error {
	// 			return gptBot.Shutdown(ctx)
	// 		},
	// 	},
	// )

	// go gptBot.Start()
	// go tgBot.Update(60)
	// go newHttpServer.Listen()

	// <-waitShutdown

}
