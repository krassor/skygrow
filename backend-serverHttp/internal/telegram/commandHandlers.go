package telegramBot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/krassor/skygrow/backend-serverHttp/internal/models/entities"
)

func (bot *Bot) list(msg *tgbotapi.MessageConfig) error {

	devices, err := bot.service.GetDevices(context.Background())
	if err != nil {
		return err
	}

	var inlineKeyboardRow []tgbotapi.InlineKeyboardButton
	var inlineNumericKeyboard tgbotapi.InlineKeyboardMarkup

	for i, device := range devices {
		buttonText := fmt.Sprintf("%s %s %s:%s", device.DeviceVendor, device.DeviceName, device.DeviceIpAddress, device.DevicePort)
		buttonId := fmt.Sprintf("%d", device.ID)

		inlineKeyboardRow = append(inlineKeyboardRow, tgbotapi.InlineKeyboardButton{Text: buttonText, CallbackData: &buttonId})

		if ((i + 1) % 2) == 0 {
			inlineNumericKeyboard.InlineKeyboard = append(inlineNumericKeyboard.InlineKeyboard, inlineKeyboardRow)
			inlineKeyboardRow = nil
		}

	}

	if inlineKeyboardRow != nil {
		inlineNumericKeyboard.InlineKeyboard = append(inlineNumericKeyboard.InlineKeyboard, inlineKeyboardRow)
	}

	msg.Text = "Select device:"
	msg.ReplyMarkup = inlineNumericKeyboard
	return nil
}

func (bot *Bot) subscribe(msg *tgbotapi.MessageConfig) error {
	subscriber, err := bot.subscriber.GetSubscriberByChatId(context.Background(), msg.ChatID)
	if err != nil {
		msg.Text = "Error finding subscriber in the DB"
		return err
	}
	if (subscriber == entities.Subscriber{}) {
		_, err = bot.subscriber.CreateNewSubscriber(context.Background(), msg.ChatID, msg.ChannelUsername)
		if err != nil {
			msg.Text = "Error creating subscriber in the DB"
			return err
		}
		msg.Text = "Succeful subscribe"
		return nil
	}

	if !subscriber.IsActive {
		_, err = bot.subscriber.UpdateSubscriberByChatId(context.Background(), subscriber, true)
		if err != nil {
			msg.Text = "Error updating subscriber status in the DB"
			return err
		}
		msg.Text = "Succeful update subscriber status"
		return nil
	}
	msg.Text = "You have aleady subscribed"
	return nil
}
