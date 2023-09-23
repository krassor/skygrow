package redisBroker

import (
	"context"
	"encoding/json"
	//"fmt"

	"github.com/krassor/skygrow/tg-gpt-bot/internal/dto"

	"github.com/rs/zerolog/log"
)

func (c *RedisClient) Subscribe(ctx context.Context, channels ...string) <-chan dto.OpenaiMsg {
	ch := make(chan dto.OpenaiMsg)

	op := "SubscribeOpenAiRequest"
	log.Info().Msgf("channel %s. %s: starting subscriber...", channels, op)

	msg := dto.OpenaiMsg{}

	sub := c.Client.Subscribe(ctx, channels...)

	go func() {
		messages := sub.Channel()
		for message := range messages {
			select {
			// case <-c.shutdownChannel:
			// 	close(ch)
			// 	return
			case <-ctx.Done():
				close(ch)
				return
			default:
				log.Info().Msgf("%s. Message: %v", op, message)
				err := json.Unmarshal([]byte(message.Payload), &msg)
				if err != nil {
					log.Error().Msgf("%s: %s", op, err)
				}
				ch <- msg
			}
		}
	}()
	return ch
}

// func (c *RedisClient) Shutdown(ctx context.Context) error {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return fmt.Errorf("force exit SubscribeOpenAiRequest: %w", ctx.Err())
// 		default:
// 			close(c.shutdownChannel)
// 			return nil
// 		}
// 	}
// }
