package redisBroker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/krassor/skygrow/tg-gpt-bot/internal/dto"

	"github.com/rs/zerolog/log"
)

func (c *RedisClient) SubscribeOpenAiRequest(ctx context.Context) (out chan dto.OpenaiMsg) {
	op := "SubscribeOpenAiRequest"
	log.Info().Msgf("%s: starting subscriber...", op)

	msg := dto.OpenaiMsg{}
	
	sub := c.Client.Subscribe(ctx, c.requestRedisChannel)

	messages := sub.Channel()
	for message := range messages {
		select {
		case <-c.shutdownChannel:
			return
		default:
			log.Info().Msgf("%s. Message: %v", op, message)
			err := json.Unmarshal([]byte(message.Payload), &msg)
			if err != nil {
				log.Error().Msgf("%s: %s", op, err)
			}
			out <- msg
		}
	}
}

func (c *RedisClient) Shutdown(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("force exit SubscribeOpenAiRequest: %w", ctx.Err())
		default:
			close(c.shutdownChannel)
			close(c.requestChan)
			close(c.responseChan)
			return nil
		}
	}
}
