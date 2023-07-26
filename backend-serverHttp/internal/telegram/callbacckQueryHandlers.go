package telegramBot

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) callbackQueryHandle(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) error {

	id, err := strconv.Atoi(callbackQuery.Data)
	if err != nil {
		callback := tgbotapi.NewCallback(callbackQuery.ID, "Internal error")
		_, errw := bot.tgbot.Request(callback)
		err = fmt.Errorf("%w", errw)
		return err
	}

	deviceEntity, err := bot.service.GetDeviceById(ctx, uint(id))
	if err != nil {
		callback := tgbotapi.NewCallback(callbackQuery.ID, "Internal error")
		_, errw := bot.tgbot.Request(callback)
		err = fmt.Errorf("%w", errw)
		return err
	}

	var status string
	if deviceEntity.DeviceStatus {
		status = "ONLINE"
	} else {
		status = "OFFLINE"
	}

	callbackData := fmt.Sprintf(
		"Device %s %s is %s",
		deviceEntity.DeviceVendor,
		deviceEntity.DeviceName,
		status,
	)
	callback := tgbotapi.NewCallback(callbackQuery.ID, callbackData)
	_, err = bot.tgbot.Request(callback)
	if err != nil {
		return err
	}

	replyMsg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, callbackData)
	_, err = bot.tgbot.Send(replyMsg)
	if err != nil {
		return err
	}
	return nil
}
