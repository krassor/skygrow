package deepseek

import (
	"context"
	"fmt"
	"log/slog"

	"app/main.go/internal/config"
	"app/main.go/internal/utils/logger/sl"

	dsmod "github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
)

type Cache interface {
	Save(ctx context.Context, userID int64, history []dsmod.ChatCompletionMessage) error
	Get(ctx context.Context, userID int64) ([]dsmod.ChatCompletionMessage, error)
	Delete(ctx context.Context, id int64) error
}

type DeepSeek struct {
	logger          *slog.Logger
	config          *config.Config
	Client          *dsmod.Client
	cache           Cache
	shutdownChannel chan struct{}
	ctx             context.Context
	cancel          context.CancelFunc
}

func NewClient(
	logger *slog.Logger,
	config *config.Config,
	cache Cache,
) *DeepSeek {

	op := "deepseek.NewClient()"
	log := logger.With(
		slog.String("op", op),
	)

	aiToken := config.BotConfig.AI.AIApiToken

	ctx, cancel := context.WithCancel(context.Background())

	client := dsmod.NewClient(aiToken)

	log.Info("Creating deepseek client")

	return &DeepSeek{
		logger:          logger,
		config:          config,
		Client:          client,
		cache:           cache,
		shutdownChannel: make(chan struct{}),
		ctx:             ctx,
		cancel:          cancel,
	}
}

func (ds *DeepSeek) ProcessMessage(ctx context.Context, userID int64, message string) (string, error) {
	op := "deepseek.ProcessMessage()"
	log := ds.logger.With(
		slog.String("op", op),
		slog.Int64("userID", userID),
	)

	// Получаем историю сообщений из cache
	history, err := ds.cache.Get(ctx, userID)
	if err != nil {
		sl.Err(err)
	}

	// Добавляем новое сообщение
	history = append(history, dsmod.ChatCompletionMessage{
		Role:    constants.ChatMessageRoleUser,
		Content: message,
	})

	// Обрезаем историю если превышен лимит сообщений
	history = ds.truncateHistory(history)

	resp, err := ds.Client.CreateChatCompletion(
		ctx,
		&dsmod.ChatCompletionRequest{
			Model:    dsmod.DeepSeekChat,
			Messages: history,
		},
	)

	if err != nil {
		return "", err
	}

	//response := fmt.Sprintf("%s\n-----------\nCompletion tokens usage: %v\nPrompt tokens usage%v\nTotal tokens usage: %v", resp.Choices[0].Message.Content, resp.Usage.CompletionTokens, resp.Usage.PromptTokens, resp.Usage.TotalTokens)

	//Save response from openai GPT bot as an assistant response
	msg := dsmod.ChatCompletionMessage{
		Role:    constants.ChatMessageRoleAssistant,
		Content: resp.Choices[0].Message.Content,
	}
	// Добавляем ответ в кэш
	history = append(history, msg)
	err = ds.cache.Save(ctx, userID, history)
	if err != nil {
		log.Error("", slog.String("error", err.Error()))
	}

	return resp.Choices[0].Message.Content, nil
}

func (ds *DeepSeek) truncateHistory(history []dsmod.ChatCompletionMessage) []dsmod.ChatCompletionMessage {

	if len(history) > 11 {
		return history[len(history)-6:]
	}

	//Восстанавливаем системный промт в [0] сообщении
	systemRoleMessage := make([]dsmod.ChatCompletionMessage, 1)
	systemRoleMessage[0] = dsmod.ChatCompletionMessage{
		Role:    constants.ChatMessageRoleSystem,
		Content: ds.config.BotConfig.AI.SystemRolePromt,
	}

	if history[0].Role != constants.ChatMessageRoleSystem {
		//Добавляем системное сообщение вперед
		history = append(systemRoleMessage, history...)
	} else {
		history[0] = systemRoleMessage[0]
	}

	return history
}

func (ds *DeepSeek) Shutdown(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("force exit AI client: %w", ctx.Err())
		default:
			ds.cancel()
			return nil
		}
	}
}
