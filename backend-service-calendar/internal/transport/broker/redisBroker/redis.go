package redisBroker

import (
	"context"
	"fmt"
	"os"

	//"github.com/krassor/skygrow/tg-gpt-bot/internal/dto"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisClient struct {
	Client               *redis.Client
	//shutdownChannel      chan struct{}
}

func NewRedisClient() *RedisClient {
	op := "NewRedisClient"
	redisHost, ok := os.LookupEnv("REDIS_HOST")
	if !ok {
		log.Error().Msgf("Cannot find REDIS_HOST env")
		return nil
	}
	redisPort, ok := os.LookupEnv("REDIS_PORT")
	if !ok {
		log.Error().Msgf("Cannot find REDIS_PORT env")
		return nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Error().Msgf("%s:%s", op, err)
	}
	log.Info().Msgf("connected to redis: %s:%s", redisHost, redisPort)

	return &RedisClient{
		Client:               rdb,
		//shutdownChannel:      make(chan struct{}),
	}
}
