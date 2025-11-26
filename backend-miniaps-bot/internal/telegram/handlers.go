package telegramBot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"log/slog"

	"app/main.go/internal/utils/logger/sl"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type sendFunction func(inputMsg *tgbotapi.Message, replyText string) error

func (bot *Bot) defaultHandler(ctx context.Context, update *tgbotapi.Update, sendFunc sendFunction) {
	op := "bot.defaultHandler"
	log := bot.log.With(
		slog.String("op", op),
	)
	log.Debug("input message",
		slog.String("user id", strconv.FormatInt(update.Message.From.ID, 10)),
		slog.String("user name", update.Message.From.UserName),
		slog.String("first name", update.Message.From.FirstName),
		slog.String("last name", update.Message.From.LastName),
		slog.String("message id", strconv.Itoa(update.Message.MessageID)),
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(bot.cfg.AIConfig.Timeout)*time.Second)
	defer cancel()

	// Этот анонимный функциональный литерал используется для отправки сообщения о наборе текста в чате.
	// Он запускается в отдельной горутине, чтобы не блокировать основной поток выполнения.
	//
	// Функция использует `select` для проверки, не был ли контекст отменен.
	// Если контекст был отменен, функция возвращает управление.
	// В противном случае, функция отправляет сообщение о наборе текста в чате.
	//
	// После отправки сообщения, функция "спит" на 2 секунды, чтобы имитировать задержку при наборе текста.
	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			bot.tgbot.Send(tgbotapi.NewChatAction(update.FromChat().ID, tgbotapi.ChatTyping))
		}
		time.Sleep(2 * time.Second)
	}()

	response, err := bot.AIBot.ProcessMessage(
		ctxTimeout,
		update.Message.From.ID,
		bot.textFilter(update.Message.Text),
	)
	cancel()

	if err != nil {
		sl.Err(err)
		log.Error("failed to process message with AI", sl.Err(err))
	}

	log.Debug("Got response from AI", slog.String("response", response))

	err = sendFunc(update.Message, response)
	if err != nil {
		log.Error("failed to send response to user", sl.Err(err))
	} else {
		log.Debug("Sent response to user")
	}

}

func (bot *Bot) stubHandler(ctx context.Context, update *tgbotapi.Update) {
	op := "bot.handleStub"
	log := bot.log.With(
		slog.String("op", op),
	)
	select {
	case <-ctx.Done():
		return
	default:
		log.Info("input message",
			slog.String("user id", strconv.FormatInt(update.Message.From.ID, 10)),
			slog.String("user name", update.Message.From.UserName),
			slog.String("first name", update.Message.From.FirstName),
			slog.String("last name", update.Message.From.LastName),
			slog.String("message id", strconv.Itoa(update.Message.MessageID)),
		)
	}

}

func (bot *Bot) commandHandler(ctx context.Context, update *tgbotapi.Update, sendFunc sendFunction) error {
	// op := "bot.commandHandle"
	// Extract the command from the Message.
	// log := bot.log.With(
	// 	slog.String("op", op),
	// )

	msg := update.Message

	switch update.Message.Command() {

	case "askbot":
		ctxTimeout, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		response, err := bot.AIBot.ProcessMessage(
			ctxTimeout,
			update.Message.From.ID,
			bot.textFilter(update.Message.Text),
		)
		if err != nil {
			sl.Err(err)
		}

		err = sendFunc(update.Message, response)
		if err != nil {
			sl.Err(err)
		}

	case "start":
		replyText := fmt.Sprintf("Hi, %s! Ask your questions.", msg.From.UserName)
		err := bot.sendMessage(msg, replyText)
		if err != nil {
			return fmt.Errorf("tgbot.commandHandle: %w", err)
		}

	case "calendar":
		err := bot.sendMenu(msg, "Calendar")
		if err != nil {
			return fmt.Errorf("tgbot.commandHandle: %w", err)
		}

	default:
		replyText := "I don't know this command"
		err := bot.sendReplyMessage(msg, replyText)
		if err != nil {
			return fmt.Errorf("tgbot.commandHandle: %w", err)
		}
	}

	return nil
}

// textFilter processes the input message by removing the "/askbot" prefix
// and filtering out words that start with "@". This function is used to
// clean up user input before processing it further.
//
// Parameters:
//   - msg: A string containing the original message text.
//
// Returns:
//
//	A string with the "/askbot" prefix removed and any words starting with "@" filtered out.
func (bot *Bot) textFilter(msg string) string {

	msgText := strings.TrimPrefix(msg, "/askbot ")

	words := strings.Split(msgText, " ")
	var filteredWords []string
	for _, word := range words {
		if !strings.HasPrefix(word, "@") {
			filteredWords = append(filteredWords, word)
		}
	}
	msgText = strings.Join(filteredWords, " ")
	return msgText
}
