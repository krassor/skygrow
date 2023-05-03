package telegramBot

import (
	"context"
	"fmt"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/krassor/skygrow/internal/openai"
	"github.com/rs/zerolog/log"
)

type Bot struct {
	tgbot  *tgbotapi.BotAPI
	gptBot *openai.GPTBot
}

func NewBot(gptBot *openai.GPTBot) *Bot {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGBOT_APITOKEN"))
	if err != nil {
		log.Error().Msgf("Error auth telegram bot: %s", err)
	}
	//TODO: add to env BOTDEBUG
	bot.Debug = false

	log.Info().Msgf("Authorized on account %s", bot.Self.UserName)

	return &Bot{
		tgbot:  bot,
		gptBot: gptBot,
	}
}

func (bot *Bot) Update(ctx context.Context, updateTimeout int) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = updateTimeout

	updates := bot.tgbot.GetUpdatesChan(updateConfig)

	//TODO: make goroutine with check update channel close
	for update := range updates {
		log.Info().Msgf("Input message: %v", update.Message)

		if update.Message == nil { // ignore any non-Message updates
			log.Warn().Msgf("tgbot warn: Not message: %v", update.Message)
			continue
		}

		// Проверяем, если сообщение адресовано самому боту
		if update.Message.Chat.ID == bot.tgbot.Self.ID {
			log.Info().Msgf("Self message: %s", update.Message.Text)
			replyText, err := bot.sendMessageToOpenAI(update.Message)
			if err != nil {
				log.Error().Msgf("Error tgbot.update: %w", err)
			}
			err = bot.sendReplyMessage(update.Message, replyText)
			if err != nil {
				log.Error().Msgf("Error tgbot.update: %w", err)
			}
			continue
		}

		// если сообщение адресовано каналу, в котором находится бот
		if update.Message.Chat.IsChannel() && update.Message.Chat.UserName == bot.tgbot.Self.UserName {
			log.Info().Msgf("Channel message from: %s", update.Message.From.UserName)
			replyText, err := bot.sendMessageToOpenAI(update.Message)
			if err != nil {
				log.Error().Msgf("Error tgbot.update: %w", err)
			}
			err = bot.sendReplyMessage(update.Message, replyText)
			if err != nil {
				log.Error().Msgf("Error tgbot.update: %w", err)
			}
			continue
		}

		// Проверяем, если сообщение является ответом на сообщение бота
		if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.ID == bot.tgbot.Self.ID {
			log.Info().Msgf("Reply message from: %s", update.Message.From.UserName)
			replyText, err := bot.sendMessageToOpenAI(update.Message)
			if err != nil {
				log.Error().Msgf("Error tgbot.update: %w", err)
			}
			err = bot.sendReplyMessage(update.Message, replyText)
			if err != nil {
				log.Error().Msgf("Error tgbot.update: %w", err)
			}
			continue
		}

		//Check if message is a command
		if update.Message.IsCommand() {
			log.Info().Msgf("tgbot.update receive command from %s: %s, text: %s", update.Message.From, update.Message.Command(), update.Message.Text)

			if err := bot.commandHandle(update.Message); err != nil {
				log.Error().Msgf("Error tgbot.update: %w", err)
			}
			continue
		}

	}
	log.Info().Msgf("exit tgbot routine")
}

func (bot *Bot) commandHandle(msg *tgbotapi.Message) error {

	// Extract the command from the Message.

	switch msg.Command() {
	case "askbot":
		replyText, err := bot.sendMessageToOpenAI(msg)
		if err != nil {
			return fmt.Errorf("Error tgbot.commandHandle: %w", err)
		}
		err = bot.sendReplyMessage(msg, replyText)
		if err != nil {
			return fmt.Errorf("Error tgbot.commandHandle: %w", err)
		}
	default:
		replyText := "I don't know this command"
		err := bot.sendReplyMessage(msg, replyText)
		if err != nil {
			return fmt.Errorf("Error tgbot.commandHandle: %w", err)
		}
	}

	return nil
}

func (bot *Bot) sendMessageToOpenAI(msg *tgbotapi.Message) (string, error) {
	msgText := strings.TrimPrefix(msg.Text, "/askbot ")

	reply, err := bot.gptBot.CreateChatCompletion(msgText)
	if err != nil {
		return "", fmt.Errorf("Error bot.sendMessageToOpenAI: %w", err)
	}
	return reply, nil
}

func (bot *Bot) sendReplyMessage(inputMsg *tgbotapi.Message, replyText string) error {
	replyMsg := tgbotapi.NewMessage(inputMsg.Chat.ID, "")
	replyMsg.ReplyToMessageID = inputMsg.MessageID
	replyMsg.Text = replyText

	_, err := bot.tgbot.Send(replyMsg)
	if err != nil {
		return fmt.Errorf("Error tgbot.sendRelyMessage: %w", err)
	}
	return nil
}

func (bot *Bot) Shutdown(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("error shutdown telegram bot: %w", ctx.Err())
		default:
			bot.tgbot.StopReceivingUpdates()
			return nil
		}
	}
}
