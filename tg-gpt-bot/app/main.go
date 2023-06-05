package main

import (
	"context"
	"time"

	"github.com/krassor/skygrow/internal/config"
	"github.com/krassor/skygrow/internal/graceful"
	"github.com/krassor/skygrow/internal/logger"
	"github.com/krassor/skygrow/internal/openai"
	"github.com/krassor/skygrow/internal/repository"
	telegramBot "github.com/krassor/skygrow/internal/telegram"
)

func main() {

	logger.InitLogger()

	config := config.InitConfig()
	repo := repository.NewMessageRepository()

	gptBot := openai.NewGPTBot(config, repo)
	tgBot := telegramBot.NewBot(config, gptBot)

	maxSecond := 15 * time.Second
	waitShutdown := graceful.GracefulShutdown(
		context.Background(),
		maxSecond,
		map[string]graceful.Operation{
			"tgBot": func(ctx context.Context) error {
				return tgBot.Shutdown(ctx)
			},
			// "tgBot": func(ctx context.Context) error {
			// 	return deviceTgBot.Shutdown(ctx)
			// },
		},
	)

	go tgBot.Update(context.Background(), 60)

	<-waitShutdown
}
