package telegramBot

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

	case "setpromtfile":
		replyText := ""
		isAdmin, err := bot.isAdmin(update.Message)

		log.Debug("setpromtfile",
			slog.String("user name", update.Message.From.UserName),
			slog.String("message", update.Message.Text),
			slog.String("is admin", strconv.FormatBool(isAdmin)),
		)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if isAdmin {
			// Проверка наличия файла в сообщении
			if update.Message.Document == nil {
				log.Info(
					"no attached file due command /setpromtfile",
					slog.String("user name", update.Message.From.UserName),
					slog.String("message", update.Message.Text),
				)
				replyText = "No file attached"
				err := sendFunc(update.Message, replyText)
				e := fmt.Errorf("No file attached")
				if err != nil {
					return fmt.Errorf("%s: %w", op, err)
				}
				return fmt.Errorf("%s: %w", op, e)
			}

			// Обработка файла
			fileID := update.Message.Document.FileID
			log.Info(
				"Received command with file",
				slog.String("user name", update.Message.From.UserName),
				slog.String("message", update.Message.Text),
				slog.String("file_id", fileID),
			)

			//Получаем file_path
			fileURL, err := bot.tgbot.GetFileDirectURL(fileID)
			if err != nil {
				replyText = "Cannot download file. PLease try again"
				e := sendFunc(update.Message, replyText)
				if e != nil {
					return fmt.Errorf("%s: %w", op, e)
				}
				return fmt.Errorf("%s: %w", op, err)
			}

			// Делаем HTTP GET-запрос по URL
			resp, err := http.Get(fileURL)
			if err != nil {
				replyText = "Cannot download file. PLease try again"
				e := sendFunc(update.Message, replyText)
				if e != nil {
					return fmt.Errorf("%s: %w", op, e)
				}
				return fmt.Errorf("%s: %w", op, err)
			}
			defer resp.Body.Close()

			//Сохраняем файл на диск
			buf := make([]byte, resp.ContentLength)

			_, err = resp.Body.Read(buf)
			if err != nil {
				replyText = "Cannot download file. PLease try again"
				e := sendFunc(update.Message, replyText)
				if e != nil {
					return fmt.Errorf("%s: %w", op, e)
				}
				return fmt.Errorf("%s: %w", op, err)
			}

			filePath := filepath.Join(bot.cfg.BotConfig.AI.PromtFilePath, bot.cfg.BotConfig.AI.PromtFileName)

			err = os.WriteFile(filePath, buf, 0775)
			if err != nil {
				replyText = "Cannot save file. PLease try again"
				e := sendFunc(update.Message, replyText)
				if e != nil {
					return fmt.Errorf("%s: %w", op, e)
				}
				return fmt.Errorf("%s: %w", op, err)
			}

			log.Info(
				"promt file saved",
				slog.String("user name", update.Message.From.UserName),
				slog.String("message", update.Message.Text),
				slog.String("file_id", fileID),
				slog.String("file_path", filePath),
			)

			//Перечитываем заново промт из файла для применения изменений
			err = bot.cfg.ReadPromtFromFile()
			if err != nil {
				replyText = "Promt file saved. But config file not updated. PLease try again"
				e := sendFunc(update.Message, replyText)
				if e != nil {
					return fmt.Errorf("%s: %w", op, e)
				}
				return fmt.Errorf("%s: %w", op, err)
			}

			log.Info(
				"Promt file saved. Config updated.",
				slog.String("user name", update.Message.From.UserName),
				slog.String("message", update.Message.Text),
				slog.String("file_id", fileID),
				slog.String("file_path", filePath),
			)

			replyText = "👍 Promt file saved. Config updated 👍"
			err = sendFunc(update.Message, replyText)
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

	case "setmodel":
		replyText := ""
		isAdmin, err := bot.isAdmin(update.Message)

		log.Debug("setmodel",
			slog.String("user name", update.Message.From.UserName),
			slog.String("message", update.Message.Text),
			slog.String("is admin", strconv.FormatBool(isAdmin)),
		)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if isAdmin {

			bot.cfg.BotConfig.AI.ModelName = strings.TrimPrefix(
				update.Message.Text, "/setmodel ")

			err = bot.cfg.Write()
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}

			log.Debug(
				"system model changed",
				slog.String("model", bot.cfg.BotConfig.AI.ModelName),
			)

			replyText = "👍 Model changed 👍"
			err := sendFunc(update.Message, replyText)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}

	case "getmodel":

		replyText := ""
		isAdmin, err := bot.isAdmin(update.Message)

		log.Debug("getmodel",
			slog.String("user name", update.Message.From.UserName),
			slog.String("message", update.Message.Text),
			slog.String("is admin", strconv.FormatBool(isAdmin)),
		)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if isAdmin {

			replyText = bot.cfg.BotConfig.AI.ModelName
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
