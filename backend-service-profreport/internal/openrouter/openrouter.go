package openrouter

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"app/main.go/internal/config"

	openrouter "github.com/revrost/go-openrouter"
)

const (
	retryCount    int           = 3
	retryDuration time.Duration = 3 * time.Second
)

type Openrouter struct {
	logger          *slog.Logger
	config          *config.Config
	Client          *openrouter.Client
	shutdownChannel chan struct{}
}

func NewClient(
	logger *slog.Logger,
	config *config.Config,
) *Openrouter {

	op := "Openrouter.NewClient()"
	log := logger.With(
		slog.String("op", op),
	)

	client := openrouter.NewClient(
		config.BotConfig.AI.AIApiToken,
	)

	log.Info("Creating deepseek client")

	return &Openrouter{
		logger:          logger,
		config:          config,
		Client:          client,
		shutdownChannel: make(chan struct{}),
	}
}

func (or *Openrouter) CreateChatCompletion(ctx context.Context, logger *slog.Logger, message string) (string, error) {
	op := "deepseek.CreateChatCompletion()"
	log := logger.With(
		slog.String("op", op),
	)

	//log.Debug("input message", slog.String("message", message))
	var resp openrouter.ChatCompletionResponse
	var err error
	for retry := range retryCount {
		var r openrouter.ChatCompletionResponse
		var e error
		select {
		case <-or.shutdownChannel:
			return "", fmt.Errorf("shutdown openrouter client")
		default:
			r, e = or.Client.CreateChatCompletion(
				ctx,
				openrouter.ChatCompletionRequest{
					Model: or.config.BotConfig.AI.ModelName,
					Messages: []openrouter.ChatCompletionMessage{
						openrouter.SystemMessage(or.config.BotConfig.AI.SystemRolePromt),
						openrouter.UserMessage(message),
					},
				},
			)
		}
		if e != nil && isRateLimitError(err) {
			err = e
			log.Error(
				"rate limit CreateChatCompletion 429",
				slog.String("error", err.Error()),
				slog.Int("retry", retry),
			)
			time.Sleep(retryDuration)
			continue
		}
		resp = r
		err = e
		break
	}

	if err != nil {
		// log.Error("error creating chat completion request", slog.String("error", err.Error()))
		return "", fmt.Errorf("error creating chat completion request: %w", err)
	}

	log.Debug("received chat completion response", slog.Any("response role", resp.Choices[0].Message.Role))

	return resp.Choices[0].Message.Content.Text, nil
}

func isRateLimitError(err error) bool {
	// Если библиотека возвращает ошибку с полем Code:
	// if e, ok := err.(interface{ Code() int }); ok && e.Code() == 429 {
	// 	return true
	// }
	// Или проверка по строке ошибки (менее надёжно):
	if err != nil {
		return strings.Contains(err.Error(), "HTTP 429")
	} else {
		return false
	}
}

func (or *Openrouter) Shutdown(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("force exit AI client: %w", ctx.Err())
		default:
			close(or.shutdownChannel)
			return nil
		}
	}
}
