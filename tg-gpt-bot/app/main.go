package main

import (
	"context"
	"time"

	"github.com/krassor/skygrow/tg-gpt-bot/internal/config"
	"github.com/krassor/skygrow/tg-gpt-bot/internal/graceful"
	"github.com/krassor/skygrow/tg-gpt-bot/internal/logger"
	"github.com/krassor/skygrow/tg-gpt-bot/internal/openai"
	"github.com/krassor/skygrow/tg-gpt-bot/internal/repository"
	telegramBot "github.com/krassor/skygrow/tg-gpt-bot/internal/telegram"
	"github.com/krassor/skygrow/tg-gpt-bot/internal/transport/httpServer"
	"github.com/krassor/skygrow/tg-gpt-bot/internal/transport/httpServer/routers"
)

func main() {
	//123

	logger.InitLogger()

	config := config.InitConfig()
	repo := repository.NewMessageRepository()

	gptBot := openai.NewGPTBot(config, repo)
	tgBot := telegramBot.NewBot(config, gptBot)

	botRouter := routers.NewBotRouter()
	httpServer := httpServer.NewHttpServer(botRouter)

	maxSecond := 15 * time.Second
	waitShutdown := graceful.GracefulShutdown(
		context.Background(),
		maxSecond,
		map[string]graceful.Operation{
			"tgBot": func(ctx context.Context) error {
				return tgBot.Shutdown(ctx)
			},
			"httpServer": func(ctx context.Context) error {
				return httpServer.Shutdown(ctx)
			},
			// "tgBot": func(ctx context.Context) error {
			// 	return deviceTgBot.Shutdown(ctx)
			// },
		},
	)

	go tgBot.Update(60)
	go httpServer.Listen()

	<-waitShutdown
}
