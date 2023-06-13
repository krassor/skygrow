package telegramBot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func (bot *Bot) replyHandler(msg *tgbotapi.Message) {
	log.Info().Msgf("Reply message from: %s", msg.From.UserName)

	replyText, err := bot.sendMessageToOpenAI(msg)
	if err != nil {
		log.Error().Msgf("Error tgbot.update: %v", err)
		return
	}
	err = bot.sendReplyMessage(msg, replyText)
	if err != nil {
		log.Error().Msgf("Error tgbot.update: %v", err)
		return
	}
}

func (bot *Bot) privateHandler(msg *tgbotapi.Message) {
	log.Info().Msgf("Self message: %s", msg.Text)

	replyText, err := bot.sendMessageToOpenAI(msg)
	if err != nil {
		log.Error().Msgf("Error tgbot.update: %v", err)
		return
	}
	err = bot.sendReplyMessage(msg, replyText)
	if err != nil {
		log.Error().Msgf("Error tgbot.update: %v", err)
		return
	}
}

func (bot *Bot) channelHandler(msg *tgbotapi.Message) {
	log.Info().Msgf("Channel: %s. Message from: %s", msg.Chat.Title, msg.From.UserName)

	replyText, err := bot.sendMessageToOpenAI(msg)
	if err != nil {
		log.Error().Msgf("Error tgbot.update: %v", err)
		return
	}
	err = bot.sendReplyMessage(msg, replyText)
	if err != nil {
		log.Error().Msgf("Error tgbot.update: %v", err)
		return
	}
}

func (bot *Bot) commandHandle(msg *tgbotapi.Message) error {

	// Extract the command from the Message.

	switch msg.Command() {
	case "setsystempromt":

		log.Info().Msgf("Input setsystempromt")

		replyText := ""
		isAdmin, err := bot.isAdmin(msg)

		log.Info().Msgf("setsystempromt. %s is admin: %v", msg.From.UserName, isAdmin)

		if err != nil {
			return fmt.Errorf("Error tgbot.commandHandle: %w", err)
		}
		if isAdmin {
			openAIConfig, err := bot.botConfig.ReadOpenAIConfig()

			log.Info().Msgf("read openAI config: %v", openAIConfig)

			if err != nil {
				return fmt.Errorf("Error tgbot.commandHandle: %w", err)
			}

			openAIConfig.SystemRolePromt = strings.TrimPrefix(msg.Text, "/setsystempromt ")

			err = bot.botConfig.WriteOpenAIConfig(&openAIConfig)
			if err != nil {
				return fmt.Errorf("Error tgbot.commandHandle: %w", err)
			}

			replyText = "👍 System role promt changed 👍"
			err = bot.sendReplyMessage(msg, replyText)
			if err != nil {
				return fmt.Errorf("Error tgbot.commandHandle: %w", err)
			}
		}

	case "getsystempromt":

		log.Info().Msgf("Input getsystempromt")

		replyText := ""
		isAdmin, err := bot.isAdmin(msg)

		log.Info().Msgf("getsystempromt. %s is admin: %v", msg.From.UserName, isAdmin)

		if err != nil {
			return fmt.Errorf("Error tgbot.commandHandle: %w", err)
		}
		if isAdmin {
			openAIConfig, err := bot.botConfig.ReadOpenAIConfig()

			log.Info().Msgf("read openAI config: %v", openAIConfig)

			if err != nil {
				return fmt.Errorf("Error tgbot.commandHandle: %w", err)
			}

			replyText = openAIConfig.SystemRolePromt

			err = bot.sendReplyMessage(msg, replyText)
			if err != nil {
				return fmt.Errorf("Error tgbot.commandHandle: %w", err)
			}
		}

	case "askbot":
		replyText, err := bot.sendMessageToOpenAI(msg)
		if err != nil {
			return fmt.Errorf("Error tgbot.commandHandle: %w", err)
		}
		err = bot.sendReplyMessage(msg, replyText)
		if err != nil {
			return fmt.Errorf("Error tgbot.commandHandle: %w", err)
		}
	case "start":
		replyText := fmt.Sprintf("Hi, %s! Ask your questions.", msg.From.UserName)
		err := bot.sendReplyMessage(msg, replyText)
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
