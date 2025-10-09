package deepseek

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"app/main.go/internal/config"
	"app/main.go/internal/utils/logger/sl"

	dsmod "github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
)

type Cache interface {
	Save(ctx context.Context, userID int64, history []any) error
	Get(ctx context.Context, userID int64) ([]any, error)
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

	client := dsmod.NewClient(aiToken, config.BotConfig.AI.BaseURL)

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

	log.Debug("input message", slog.String("message", message))

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

	// Convert history to []dsmod.ChatCompletionMessage
	messages := make([]dsmod.ChatCompletionMessage, len(history))
	for i, msg := range history {
		if chatMsg, ok := msg.(dsmod.ChatCompletionMessage); ok {
			messages[i] = chatMsg
		}
	}

	log.Debug("created chat completion request", slog.Any("messages", messages))

	resp, err := ds.Client.CreateChatCompletion(
		ctx,
		&dsmod.ChatCompletionRequest{
			Model:    ds.config.BotConfig.AI.ModelName,
			Messages: messages,
		},
	)

	if err != nil {
		log.Error("error creating chat completion request", slog.String("error", err.Error()))

		if isRateLimitError(err) {
			return "Provider rate limit. Please try again later.", nil
		}

		return "", fmt.Errorf("%w", err)
	}

	log.Debug("received chat completion response", slog.Any("response", resp))

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

// truncateHistory обрезает историю чата до максимально допустимой длины,
// сохраняя первое системное сообщение и последние N сообщений.
//
// Параметры:
//   - history: история сообщений в виде среза интерфейсов `any`.
//
// Возвращает:
//   - Обновлённую историю сообщений:
//     1. Гарантирует наличие системного сообщения в начале (если его нет, добавляет).
//     2. Ограничивает длину истории значением `maxHistoryLength` (30):
//        - Сохраняет первое сообщение (системное).
//        - Сохраняет последние `keepLastN` (25) сообщений.
//        - Итоговая длина: 1 (системное) + 25 = 26 сообщений.
//
// Пример:
//   Вход: [systemMsg, msg1, msg2, ..., msg40]
//   Выход: [systemMsg, msg16, msg17, ..., msg40] (сохранено 25 последних после systemMsg)
func (ds *DeepSeek) truncateHistory(history []any) []any {
	const (
		maxHistoryLength = 30
		keepLastN        = 25
	)

	// Добавление системного промпта
	if len(history) == 0 || getMessageRole(history[0]) != constants.ChatMessageRoleSystem {
		systemMsg := createSystemMessage(ds.config.BotConfig.AI.SystemRolePromt)
		history = prependSystemMessage(history, systemMsg)
	}

	// Обрезка истории
	if len(history) > maxHistoryLength {
		keepFrom := max(len(history)-keepLastN, 1)
		return append([]any{history[0]}, history[keepFrom:]...)
	}

	return history
}

// getMessageRole извлекает роль сообщения из элемента истории.
// Предполагается, что элемент содержит поле "role" с типом string.
// Реализация зависит от внутренней структуры сообщения.
func getMessageRole(msg any) string {
	if m, ok := msg.(dsmod.ChatCompletionMessage); ok {
		return m.Role
	}
	return ""
}

// createSystemMessage создает системное сообщение для чата с заданным текстом подсказки.
//
// Параметры:
//   - prompt: текст подсказки, который будет использоваться в качестве содержимого системного сообщения.
//
// Возвращает:
//   - Объект типа dsmod.ChatCompletionMessage с ролью `constants.ChatMessageRoleSystem` и указанным содержимым.
//
// Использование:
//   - Формирование начального системного сообщения для инициализации контекста чата.
//   - Добавление правил или инструкций для модели обработки диалога.
func createSystemMessage(prompt string) any {
	return dsmod.ChatCompletionMessage{
		Role:    constants.ChatMessageRoleSystem,
		Content: prompt,
	}
}

// prependSystemMessage добавляет системное сообщение в начало истории чата,
// если оно ещё не присутствует.
//
// Параметры:
//   - history: история сообщений в виде среза интерфейсов `any`.
//   - systemMsg: системное сообщение, которое нужно добавить.
//
// Возвращает:
//   - Обновлённую историю сообщений:
//     - Если первое сообщение в `history` уже является системным — возвращает `history` без изменений.
//     - В противном случае — возвращает новый срез, где `systemMsg` добавлен в начало `history`.
func prependSystemMessage(history []any, systemMsg any) []any {
	if len(history) > 0 && getMessageRole(history[0]) == constants.ChatMessageRoleSystem {
		return history
	}
	return append([]any{systemMsg}, history...)
}

func isRateLimitError(err error) bool {
	// Если библиотека возвращает ошибку с полем Code:
	// if e, ok := err.(interface{ Code() int }); ok && e.Code() == 429 {
	// 	return true
	// }
	// Или проверка по строке ошибки (менее надёжно):
	return strings.Contains(err.Error(), "HTTP 429")
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
