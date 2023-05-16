package telegramBot

import (
	"context"
	"fmt"
	"os"
	"strings"
	"unicode/utf16"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/krassor/skygrow/internal/config"
	"github.com/krassor/skygrow/internal/openai"
	"github.com/rs/zerolog/log"
)

type Bot struct {
	tgbot     *tgbotapi.BotAPI
	gptBot    *openai.GPTBot
	botConfig *config.AppConfig
}

func NewBot(botConfig *config.AppConfig, gptBot *openai.GPTBot) *Bot {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGBOT_APITOKEN"))
	if err != nil {
		log.Error().Msgf("Error auth telegram bot: %s", err)
	}
	//TODO: add to env BOTDEBUG
	bot.Debug = false

	log.Info().Msgf("Authorized on account %s", bot.Self.UserName)

	return &Bot{
		tgbot:     bot,
		gptBot:    gptBot,
		botConfig: botConfig,
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

		//Check if message is a command
		if update.Message.IsCommand() {
			log.Info().Msgf("tgbot.update receive command from %s: %s, text: %s", update.Message.From, update.Message.Command(), update.Message.Text)

			if err := bot.commandHandle(update.Message); err != nil {
				log.Error().Msgf("Error tgbot.update: %v", err)
			}
			continue
		}

		// Проверяем, если сообщение адресовано самому боту
		if update.Message.Chat.IsPrivate() {
			bot.privateHandler(update.Message)
			continue
		}

		// если сообщение адресовано каналу, в котором находится бот
		if (update.Message.Chat.IsChannel() || update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup()) && bot.checkBotMention(update.Message) {
			bot.channelHandler(update.Message)
			continue
		}

		// Проверяем, если сообщение является ответом на сообщение бота
		if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.ID == bot.tgbot.Self.ID {
			bot.replyHandler(update.Message)
			continue
		}

	}
	log.Info().Msgf("exit tgbot routine")
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

func (bot *Bot) sendMessageToOpenAI(msg *tgbotapi.Message) (string, error) {
	msgText := strings.TrimPrefix(msg.Text, "/askbot ")

	words := strings.Split(msgText, " ")
	var filteredWords []string
	for _, word := range words {
		if !strings.HasPrefix(word, "@") {
			filteredWords = append(filteredWords, word)
		}
	}
	msgText = strings.Join(filteredWords, " ")

	reply, err := bot.gptBot.CreateChatCompletion(msg.From.UserName, msgText)
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
		return fmt.Errorf("Error tgbot.sendReplyMessage: %w", err)
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
