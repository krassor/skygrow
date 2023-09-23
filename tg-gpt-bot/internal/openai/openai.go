package openai

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/krassor/skygrow/tg-gpt-bot/internal/config"
	"github.com/krassor/skygrow/tg-gpt-bot/internal/dto"
	"github.com/krassor/skygrow/tg-gpt-bot/internal/repository"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIMsgBroker interface {
	Publish(ctx context.Context, channel string, msg dto.OpenaiMsg)
	Subscribe(ctx context.Context, channels ...string) <-chan dto.OpenaiMsg
}

type GPTBot struct {
	openAIClient    *openai.Client
	repo            repository.MessageRepository
	botConfig       *config.AppConfig
	broker          OpenAIMsgBroker
	shutdownChannel chan struct{}
	ctx             context.Context
	cancel          context.CancelFunc
}

const (
	brokerChannelSub string = "openai.request"
	brokerChannelPub string = "openai.response"
)

func NewGPTBot(botConfig *config.AppConfig, repo repository.MessageRepository, broker OpenAIMsgBroker) *GPTBot {
	openAiToken, ok := os.LookupEnv("OPENAI_TOKEN")
	if !ok {
		log.Error().Msgf("Cannot find openai token env")
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())

	client := openai.NewClient(openAiToken)
	log.Info().Msgf("Connecting to openAI")
	return &GPTBot{
		openAIClient:    client,
		repo:            repo,
		botConfig:       botConfig,
		broker:          broker,
		shutdownChannel: make(chan struct{}),
		ctx:             ctx,
		cancel:          cancel,
	}
}

func (GPTBot *GPTBot) CreateChatCompletion(ctx context.Context, openaiMsg dto.OpenaiMsg) {
	op := "GPTBot.CreateChatCompletion"
	log.Info().Msgf("GPTBot get input message: %s", openaiMsg.Msg)

	msg := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: openaiMsg.Msg,
	}

	err := GPTBot.repo.SaveUserMessage(ctx, openaiMsg.UserId, msg)
	if err != nil {
		log.Error().Msgf("%s: %v", op, err)
	}

	messages, err := GPTBot.repo.LoadUserMessages(ctx, openaiMsg.UserId)
	if err != nil {
		log.Error().Msgf("%s: %v", op, err)
	}

	openAIConfig, err := GPTBot.botConfig.ReadOpenAIConfig()
	if err != nil {
		log.Error().Msgf("%s: %v", op, err)
	}

	log.Info().Msgf("%s. read openAIConfig: %v", op, openAIConfig)

	systemRoleMessage := make([]openai.ChatCompletionMessage, 1)
	systemRoleMessage[0] = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: openAIConfig.SystemRolePromt,
	}
	//Добавляем системное сообщение вперед
	messages = append(systemRoleMessage, messages...)
	log.Info().Msgf("%s. promt messages: %v", op, messages)

	resp, err := GPTBot.openAIClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo0613,
			Temperature: openAIConfig.Temperature,
			N:           openAIConfig.N,
			MaxTokens:   openAIConfig.MaxTokens,
			Messages:    messages,
		},
	)

	if err != nil {
		log.Error().Msgf("%s: %v", op, err)
	}

	//response := fmt.Sprintf("%s\n-----------\nCompletion tokens usage: %v\nPromt tokens usage%v\nTotal tokens usage: %v", resp.Choices[0].Message.Content, resp.Usage.CompletionTokens, resp.Usage.PromptTokens, resp.Usage.TotalTokens)

	//Save response from openai GPT bot as an assistent response
	msg = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: resp.Choices[0].Message.Content,
	}

	err = GPTBot.repo.SaveUserMessage(context.Background(), openaiMsg.UserId, msg)
	if err != nil {
		log.Error().Msgf("%s: %v", op, err)
	}

	if resp.Usage.TotalTokens > 3072 {
		_, err := GPTBot.repo.DeleteFirstPromt(context.Background(), openaiMsg.UserId)
		if err != nil {
			log.Error().Msgf("%s: %v", op, err)
		}
		_, err = GPTBot.repo.DeleteFirstPromt(context.Background(), openaiMsg.UserId)
		if err != nil {
			log.Error().Msgf("%s: %v", op, err)
		}
		//log.Printf("Deleted first promt: %v", del)
	}

	if resp.Usage.TotalTokens > 2048 {
		_, err := GPTBot.repo.DeleteFirstPromt(context.Background(), openaiMsg.UserId)
		if err != nil {
			log.Error().Msgf("%s: %v", op, err)
		}
		//log.Printf("Deleted first promt: %v", del)
	}

	observeTotalTokensUsage(resp.Usage.TotalTokens, openaiMsg.UserId)

	GPTBot.broker.Publish(ctx, brokerChannelPub, dto.OpenaiMsg{
		Source: openaiMsg.Source,
		ChatId: openaiMsg.ChatId,
		UserId: openaiMsg.UserId,
		MsgId: openaiMsg.MsgId,
		Msg:    resp.Choices[0].Message.Content,
	})
}

func (GPTBot *GPTBot) Start() {
	op := "GPTBot Start()"

	updates := GPTBot.broker.Subscribe(GPTBot.ctx, brokerChannelSub)
	for update := range updates {

		log.Info().Msgf("%s: %v", op, update)
		go GPTBot.CreateChatCompletion(GPTBot.ctx, update)

	}

}

func (GPTBot *GPTBot) Shutdown(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("force exit GPTBot: %w", ctx.Err())
		default:
			GPTBot.cancel()
			return nil
		}
	}
}
