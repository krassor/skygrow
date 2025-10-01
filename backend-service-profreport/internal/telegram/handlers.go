package telegramBot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type sendFunction func(inputMsg *tgbotapi.Message, replyText string) error

func (bot *Bot) commandHandler(ctx context.Context, update *tgbotapi.Update, sendFunc sendFunction) error {
	op := "bot.commandHandle"
	// Extract the command from the Message.
	log := bot.log.With(
		slog.String("op", op),
	)

	msg := update.Message

	switch update.Message.Command() {
	case "setsystempromt":
		replyText := ""
		isAdmin, err := bot.isAdmin(update.Message)

		log.Debug("setsystempromt",
			slog.String("user name", update.Message.From.UserName),
			slog.String("message", update.Message.Text),
			slog.String("is admin", strconv.FormatBool(isAdmin)),
		)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if isAdmin {

			bot.cfg.BotConfig.AI.SystemRolePromt = strings.TrimPrefix(
				update.Message.Text, "/setsystempromt ")

			err = bot.cfg.Write()
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}

			log.Debug(
				"system promt changed",
				slog.String("promt", bot.cfg.BotConfig.AI.SystemRolePromt),
			)

			replyText = "👍 System role promt changed 👍"
			err := sendFunc(update.Message, replyText)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}

	case "getsystempromt":

		replyText := ""
		isAdmin, err := bot.isAdmin(update.Message)

		log.Debug("getsystempromt",
			slog.String("user name", update.Message.From.UserName),
			slog.String("message", update.Message.Text),
			slog.String("is admin", strconv.FormatBool(isAdmin)),
		)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if isAdmin {

			replyText = bot.cfg.BotConfig.AI.SystemRolePromt
			err := sendFunc(update.Message, replyText)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}

	case "start":
		replyText := fmt.Sprintf("Hi, %s! Ask your questions.", msg.From.UserName)
		err := sendFunc(update.Message, replyText)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

	default:
		replyText := "I don't know this command"
		err := bot.sendReplyMessage(msg, replyText)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

// // textFilter processes the input message by removing the "/askbot" prefix
// // and filtering out words that start with "@". This function is used to
// // clean up user input before processing it further.
// //
// // Parameters:
// //   - msg: A string containing the original message text.
// //
// // Returns:
// //
// //	A string with the "/askbot" prefix removed and any words starting with "@" filtered out.
// func (bot *Bot) textFilter(msg string) string {

// 	msgText := strings.TrimPrefix(msg, "/askbot ")

// 	words := strings.Split(msgText, " ")
// 	var filteredWords []string
// 	for _, word := range words {
// 		if !strings.HasPrefix(word, "@") {
// 			filteredWords = append(filteredWords, word)
// 		}
// 	}
// 	msgText = strings.Join(filteredWords, " ")
// 	return msgText
// }
