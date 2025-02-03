package redisBroker

// import (
// 	"context"
// 	"fmt"

// 	"encoding/json"

// 	"github.com/krassor/skygrow/backend-miniaps-bot/internal/dto"
// 	"github.com/rs/zerolog/log"
// )

// func (c *RedisClient) Publish(ctx context.Context, channel string, msg dto.OpenaiMsg) error {
// 	op := "redisClient.Publish"
// 	log.Info().Msgf("%s. Channel: %s. Input message: %v,", op, channel, msg)
// 	//err := c.Client.Publish(ctx, channel, msg).Err()

// 	msgByte, err := json.Marshal(msg)
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	intCmd := c.Client.Publish(ctx, channel, msgByte)
// 	result, err1 := intCmd.Result()
// 	err = intCmd.Err()

// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	log.Info().Msgf("%s. Result: %v, Error: %v. String: \t%s", op, result, err1, intCmd.String())

// 	return nil
// }
