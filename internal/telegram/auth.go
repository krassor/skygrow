package telegramBot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) isAdmin(msg *tgbotapi.Message) (bool, error) {
	result := false

	botConfig, err := bot.botConfig.ReadBotConfig()
	if err != nil {
		return result, fmt.Errorf("Error isADmin(): %w", err)
	}

	for _, admin := range botConfig.Admins {
		if admin == msg.From.UserName {
			result = true
			break
		}
	}

	return result, err
}
