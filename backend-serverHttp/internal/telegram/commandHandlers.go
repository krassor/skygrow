package telegramBot

import (
	"context"
	"fmt"
	"unicode/utf16"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/krassor/skygrow/backend-serverHttp/internal/models/entities"
	"github.com/rs/zerolog/log"
)

func (bot *Bot) replyHandler(msg *tgbotapi.Message) {
	// log.Info().Msgf("Reply message from: %s", msg.From.UserName)

	// replyText, err := bot.sendMessageToOpenAI(msg)
	// if err != nil {
	// 	log.Error().Msgf("Error tgbot.update: %v", err)
	// 	return
	// }

	// log.Info().Msgf("Last GPT message: %s", replyText)

	// err = bot.sendReplyMessage(msg, replyText)
	// if err != nil {
	// 	log.Error().Msgf("Error tgbot.update: %v", err)
	// 	return
	// }
}

func (bot *Bot) privateHandler(msg *tgbotapi.Message) {
	// log.Info().Msgf("Self message: %sfrom: %v %s %s %s,",  msg.Text, msg.From.ID, msg.From.UserName, msg.From.LastName, msg.From.FirstName)

	// replyText, err := bot.sendMessageToOpenAI(msg)
	// if err != nil {
	// 	log.Error().Msgf("Error tgbot.update: %v", err)
	// 	return
	// }

	// log.Info().Msgf("Last GPT message: %s", replyText)

	// err = bot.sendReplyMessage(msg, replyText)
	// if err != nil {
	// 	log.Error().Msgf("Error tgbot.update: %v", err)
	// 	return
	// }
}

func (bot *Bot) channelHandler(msg *tgbotapi.Message) {
	// log.Info().Msgf("Channel: %s. Message from: %s", msg.Chat.Title, msg.From.UserName)

	// replyText, err := bot.sendMessageToOpenAI(msg)
	// if err != nil {
	// 	log.Error().Msgf("Error tgbot.update: %v", err)
	// 	return
	// }

	// log.Info().Msgf("Last GPT message: %s", replyText)

	// err = bot.sendReplyMessage(msg, replyText)
	// if err != nil {
	// 	log.Error().Msgf("Error tgbot.update: %v", err)
	// 	return
	// }
}

func (bot *Bot) commandHandle(msg *tgbotapi.Message) error {

	// Extract the command from the Message.

	switch msg.Command() {
	case "help":
		replyText := "I understand command /subscribe"
		err := bot.sendReplyMessage(msg, replyText)
		if err != nil {
			return err
		}
	case "start":
		replyText := fmt.Sprintf("Hello, %s! I'm bookorder notify bot.", msg.Chat.UserName)
		err := bot.sendReplyMessage(msg, replyText)
		if err != nil {
			return err
		}
	// case "list":
	// 	err := bot.list(&replyMsg)
	// 	if err != nil {
	// 		return err
	// 	}
	case "subscribe":
		replyText, err := bot.subscribe(msg)
		if err != nil {
			return err
		}
		err = bot.sendReplyMessage(msg, replyText)
		if err != nil {
			return err
		}
	default:
		replyText := "I don't know this command"
		err := bot.sendReplyMessage(msg, replyText)
		if err != nil {
			return err
		}
	}

	return nil
}

func (bot *Bot) checkBotMention(msg *tgbotapi.Message) bool {
	result := false
	for _, entity := range msg.Entities {
		// Проверяем тип упоминания - если это упоминание, то
		// получаем само упоминание и обрабатываем его
		if entity.Type == "mention" {
			// Encode it into utf16
			utf16EncodedString := utf16.Encode([]rune(msg.Text))
			// Decode just the piece of string I need
			runeString := utf16.Decode(utf16EncodedString[entity.Offset+1 : entity.Offset+entity.Length])
			// Transform []rune into string
			mention := string(runeString)

			log.Info().Msgf("checkMention: %v, Bot name: %v", mention, bot.tgbot.Self.UserName)
			if mention == bot.tgbot.Self.UserName {
				log.Info().Msgf("Mentioned user: %s", mention)
				result = true
				break
			}

		}
	}
	return result
}

func (bot *Bot) sendReplyMessage(inputMsg *tgbotapi.Message, replyText string) error {
	replyMsg := tgbotapi.NewMessage(inputMsg.Chat.ID, "")
	replyMsg.ReplyToMessageID = inputMsg.MessageID
	replyMsg.Text = replyText

	_, err := bot.tgbot.Send(replyMsg)
	if err != nil {
		return fmt.Errorf("Error tgbot.sendReplyMessage: %w", err)
	}
	return nil
}

// func (bot *Bot) list(msg *tgbotapi.MessageConfig) error {

// 	devices, err := bot.service.GetDevices(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	var inlineKeyboardRow []tgbotapi.InlineKeyboardButton
// 	var inlineNumericKeyboard tgbotapi.InlineKeyboardMarkup

// 	for i, device := range devices {
// 		buttonText := fmt.Sprintf("%s %s %s:%s", device.DeviceVendor, device.DeviceName, device.DeviceIpAddress, device.DevicePort)
// 		buttonId := fmt.Sprintf("%d", device.ID)

// 		inlineKeyboardRow = append(inlineKeyboardRow, tgbotapi.InlineKeyboardButton{Text: buttonText, CallbackData: &buttonId})

// 		if ((i + 1) % 2) == 0 {
// 			inlineNumericKeyboard.InlineKeyboard = append(inlineNumericKeyboard.InlineKeyboard, inlineKeyboardRow)
// 			inlineKeyboardRow = nil
// 		}

// 	}

// 	if inlineKeyboardRow != nil {
// 		inlineNumericKeyboard.InlineKeyboard = append(inlineNumericKeyboard.InlineKeyboard, inlineKeyboardRow)
// 	}

// 	msg.Text = "Select device:"
// 	msg.ReplyMarkup = inlineNumericKeyboard
// 	return nil
// }

func (bot *Bot) subscribe(msg *tgbotapi.Message) (replyText string, err error) {
	subscriber, err := bot.subscriber.GetSubscriberByChatId(context.Background(), msg.SenderChat.ID)
	if err != nil {
		replyText = "Error finding subscriber in the DB"
		return replyText, err
	}
	if (subscriber == entities.Subscriber{}) {
		_, err = bot.subscriber.CreateNewSubscriber(context.Background(), msg.SenderChat.ID, msg.From.UserName)
		if err != nil {
			replyText = "Error creating subscriber in the DB"
			return replyText, err
		}
		replyText = "Succeful subscribe"
		return replyText, nil
	}

	if !subscriber.IsActive {
		_, err = bot.subscriber.UpdateSubscriberByChatId(context.Background(), subscriber, true)
		if err != nil {
			replyText = "Error updating subscriber status in the DB"
			return replyText, err
		}
		replyText = "Succeful update subscriber status"
		return replyText, nil
	}
	replyText = "You have aleady subscribed"
	return replyText, nil
}
