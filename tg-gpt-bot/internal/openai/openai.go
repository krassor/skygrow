package openai

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/krassor/skygrow/tg-gpt-bot/internal/config"
	"github.com/krassor/skygrow/tg-gpt-bot/internal/repository"

	openai "github.com/sashabaranov/go-openai"
)

// type RepoMessages interface {
// 	SaveUserMessage(ctx context.Context, username string, message openai.ChatCompletionMessage) error
// 	LoadUserMessages(ctx context.Context, username string) ([]openai.ChatCompletionMessage, error)
// 	IsUserExist(ctx context.Context, username string) (bool, error)
// }

type GPTBot struct {
	openAIClient *openai.Client
	repo         repository.MessageRepository
	botConfig    *config.AppConfig
}

func NewGPTBot(botConfig *config.AppConfig, repo repository.MessageRepository) *GPTBot {
	openAiToken, ok := os.LookupEnv("OPENAI_TOKEN")
	if !ok {
		log.Error().Msgf("Cannot find openai token env")
		return nil
	}

	client := openai.NewClient(openAiToken)
	log.Info().Msgf("Connecting to openAI")
	return &GPTBot{
		openAIClient: client,
		repo:         repo,
		botConfig:    botConfig,
	}
}

func (GPTBot *GPTBot) CreateChatCompletion(username string, gptInput string) (string, error) {

	log.Info().Msgf("GPTBot get input message: %s", gptInput)

	msg := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: gptInput,
	}

	err := GPTBot.repo.SaveUserMessage(context.Background(), username, msg)
	if err != nil {
		return "", fmt.Errorf("openai.CreateChatCompletion error: %w", err)
	}

	messages, err := GPTBot.repo.LoadUserMessages(context.Background(), username)
	if err != nil {
		return "", fmt.Errorf("openai.CreateChatCompletion error: %w", err)
	}

	openAIConfig, err := GPTBot.botConfig.ReadOpenAIConfig()
	if err != nil {
		return "", fmt.Errorf("openai.CreateChatCompletion error: %w", err)
	}

	log.Info().Msgf("CreateChatCompletion(): read openAIConfig: %v", openAIConfig)

	systemRoleMessage := make([]openai.ChatCompletionMessage, 1)
	systemRoleMessage[0] = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: openAIConfig.SystemRolePromt,
	}
	//Добавляем системное сообщение вперед
	messages = append(systemRoleMessage, messages...)
	log.Info().Msgf("CreateChatCompletion(): promt messages: %v", messages)

	resp, err := GPTBot.openAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo0613,
			Temperature: openAIConfig.Temperature,
			N:           openAIConfig.N,
			MaxTokens:   openAIConfig.MaxTokens,
			Messages:    messages,
		},
	)

	if err != nil {
		return "", fmt.Errorf("Error GPTBot.CreateChatCompletion: %w", err)
	}

	//response := fmt.Sprintf("%s\n-----------\nCompletion tokens usage: %v\nPromt tokens usage%v\nTotal tokens usage: %v", resp.Choices[0].Message.Content, resp.Usage.CompletionTokens, resp.Usage.PromptTokens, resp.Usage.TotalTokens)

	//Save response from openai GPT bot as an assistent response
	msg = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: resp.Choices[0].Message.Content,
	}

	err = GPTBot.repo.SaveUserMessage(context.Background(), username, msg)
	if err != nil {
		return "", fmt.Errorf("openai.CreateChatCompletion error: %w", err)
	}

	if resp.Usage.TotalTokens > 3072 {
		_, err := GPTBot.repo.DeleteFirstPromt(context.Background(), username)
		if err != nil {
			return "", fmt.Errorf("openai.CreateChatCompletion error: %w", err)
		}
		_, err = GPTBot.repo.DeleteFirstPromt(context.Background(), username)
		if err != nil {
			return "", fmt.Errorf("openai.CreateChatCompletion error: %w", err)
		}
		//log.Printf("Deleted first promt: %v", del)
	}

	if resp.Usage.TotalTokens > 2048 {
		_, err := GPTBot.repo.DeleteFirstPromt(context.Background(), username)
		if err != nil {
			return "", fmt.Errorf("openai.CreateChatCompletion error: %w", err)
		}
		//log.Printf("Deleted first promt: %v", del)
	}

	observeTotalTokensUsage(resp.Usage.TotalTokens, username)

	return resp.Choices[0].Message.Content, nil
}
