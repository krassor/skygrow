package redisBroker

import (
	"context"
	"fmt"

	"github.com/krassor/skygrow/tg-gpt-bot/internal/dto"
)

func (c *RedisClient) PublishOpenaiResponse(ctx context.Context, msg dto.OpenaiMsg) error {
	op := "PublishOpenaiAnswer"

	err := c.Client.Publish(ctx, c.requestChannel, msg).Err()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
