package repository

import (
	"context"
	"os"

	"github.com/krassor/skygrow/backend-miniaps-bot/internal/repository/inMemory"
	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
)

type MessageRepository interface {
	SaveUserMessage(ctx context.Context, username string, message openai.ChatCompletionMessage) error
	LoadUserMessages(ctx context.Context, username string) ([]openai.ChatCompletionMessage, error)
	DeleteFirstPromt(ctx context.Context, username string) ([]openai.ChatCompletionMessage, error)
	IsUserExist(ctx context.Context, username string) (bool, error)
}

type User struct {
	ID          int64
	Username    string
	Lastname    string
	Firstname   string
	PhoneNumber string
}

func NewMessageRepository() MessageRepository {
	db_type, ok := os.LookupEnv("USER_MESSAGE_DB_TYPE")
	if !ok {
		log.Error().Msgf("Cannot find USER_MESSAGE_DB_TYPE env")
		return nil
	}

	switch db_type {
	case "Inmemory":
		r := inMemory.NewInMemoryRepository()
		return r
	default:
		log.Error().Msgf("USER_MESSAGE_DB_TYPE env error: Unsupported database type")
		return nil
	}
}
