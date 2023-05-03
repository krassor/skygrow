package openai

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	openai "github.com/sashabaranov/go-openai"
)

type GPTBot struct {
	openAIClient *openai.Client
}

func NewGPTBot() *GPTBot {
	openAiToken, ok := os.LookupEnv("OPENAI_TOKEN")
	if !ok {
		log.Error().Msgf("Cannot find openai token env")
		return nil
	}

	client := openai.NewClient(openAiToken)
	return &GPTBot{
		openAIClient: client,
	}
}

func (GPTBot *GPTBot) CreateChatCompletion(gptInput string) (string, error) {
	log.Info().Msgf("GPTBot get input message: %s", gptInput)
	resp, err := GPTBot.openAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Temperature: 0.5,
			N:           1,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "ты опытный специалист по выращиванию марихуаны и конопли. Отвечай только на вопросы по выращиванию растений, оборудованию для выращивания. Отвечай как можно конкретнее. На остальные вопросы отвечай, что ты только гровер и не знаешь ответы на другие вопросы",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: gptInput,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("Error GPTBot.CreateChatCompletion: %w", err)
	}

	response := fmt.Sprintf("%s\n-----------\nCompletion tokens usage: %v\nPromt tokens usage%vTotal tokens usage\n%v", resp.Choices[0].Message.Content, resp.Usage.CompletionTokens, resp.Usage.PromptTokens, resp.Usage.TotalTokens)
	return response, nil
}
