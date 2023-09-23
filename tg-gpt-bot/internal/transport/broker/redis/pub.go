package redisBroker

import (
	"context"
	"fmt"

	"github.com/krassor/skygrow/tg-gpt-bot/internal/dto"
)

func (c *RedisClient) Publish(ctx context.Context, channel string, msg dto.OpenaiMsg) error {
	op := "PublishOpenaiAnswer"

	err := c.Client.Publish(ctx, channel, msg).Err()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
