package telegramBot

import (
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
