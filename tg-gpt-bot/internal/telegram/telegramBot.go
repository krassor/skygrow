package telegramBot

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/krassor/skygrow/tg-gpt-bot/internal/config"
	"github.com/krassor/skygrow/tg-gpt-bot/internal/dto"
	"github.com/krassor/skygrow/tg-gpt-bot/internal/openai"
	"github.com/rs/zerolog/log"
)

type OpenAIMsgBroker interface {
	Publish(ctx context.Context, channel string, msg dto.OpenaiMsg)
	Subscribe(ctx context.Context, channels ...string) <-chan dto.OpenaiMsg
}

type Bot struct {
	tgbot           *tgbotapi.BotAPI
	broker          OpenAIMsgBroker
	botConfig       *config.AppConfig
	shutdownChannel chan struct{}
	ctx             context.Context
	cancel          context.CancelFunc
}

const (
	brokerChannelSub string = "openai.response"
	brokerChannelPub string = "openai.request"
	dtoSource        string = "telegram"
)

func NewBot(botConfig *config.AppConfig, gptBot *openai.GPTBot, broker OpenAIMsgBroker) *Bot {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGBOT_APITOKEN"))
	if err != nil {
		log.Error().Msgf("Error auth telegram bot: %s", err)
	}
	//TODO: add to env BOTDEBUG
	bot.Debug = false

	log.Info().Msgf("Authorized on account %s", bot.Self.UserName)

	ctx, cancel := context.WithCancel(context.Background())

	return &Bot{
		tgbot:           bot,
		broker:          broker,
		botConfig:       botConfig,
		shutdownChannel: make(chan struct{}),
		ctx:             ctx,
		cancel:          cancel,
	}
}

func (bot *Bot) Update(updateTimeout int) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = updateTimeout

	_, err := bot.tgbot.MakeRequest("deleteWebhook", tgbotapi.Params{"drop_pending_updates": "false"})
	if err != nil {
		log.Error().Msgf("bot.Update() error: cannot delete WebHook: %v", err)
	}

	updates := bot.tgbot.GetUpdatesChan(updateConfig)

	for update := range updates {
		log.Info().Msgf("Input message: %v\n", update.Message)

		if update.Message == nil { // ignore any non-Message updates
			log.Warn().Msgf("tgbot warn: Not message: %v", update.Message)
			continue
		}

		go bot.processingMessages(update)
	}
	//log.Info().Msgf("exit tgbot routine")

}

func (bot *Bot) Subscribe() {
	op := "tgBot subscribe()"

	updates := bot.broker.Subscribe(bot.ctx, brokerChannelSub)
	for update := range updates {

		log.Info().Msgf("%s: %v", op, update)
		//TODO check source and send reply to handler
		if update.Source == dtoSource {

			tgBotApiMsg := tgbotapi.Message{}

			chatId, err := strconv.Atoi(update.ChatId)
			if err != nil {
				log.Error().Msgf("%s: %v", op, err)
			}
			msgId, err := strconv.Atoi(update.MsgId)
			if err != nil {
				log.Error().Msgf("%s: %v", op, err)
			}

			tgBotApiMsg.Chat.ID = int64(chatId)
			tgBotApiMsg.MessageID = msgId
			
			go func(update dto.OpenaiMsg) {
				err = bot.sendReplyMessage(&tgBotApiMsg, update.Msg)
				if err != nil {
					log.Error().Msgf("%s: %v", op, err)
				}
			}(update)
		}

	}

}

func (bot *Bot) processingMessages(update tgbotapi.Update) {

	log.Info().Msgf("\n\t\tEnter goroutine processingMessages(), id: %v, user: %s, name: %s %s", update.Message.From.ID, update.Message.From.UserName, update.Message.From.LastName, update.Message.From.FirstName)

	select {
	case <-bot.shutdownChannel:
		return
	default:
		//Check if message is a command
		if update.Message.IsCommand() {
			log.Info().Msgf("tgbot.update receive command from %s: %s, text: %s", update.Message.From, update.Message.Command(), update.Message.Text)

			if err := bot.commandHandle(update.Message); err != nil {
				log.Error().Msgf("Error tgbot.update: %v", err)
			}
		} else

		// Проверяем, если сообщение адресовано самому боту
		if update.Message.Chat.IsPrivate() {
			bot.defaultHandler(update.Message)
		} else

		// если сообщение адресовано каналу, в котором находится бот
		if (update.Message.Chat.IsChannel() || update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup()) && bot.checkBotMention(update.Message) {
			bot.defaultHandler(update.Message)
		} else

		// Проверяем, если сообщение является ответом на сообщение бота
		if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.ID == bot.tgbot.Self.ID {
			bot.defaultHandler(update.Message)
		} else {
			log.Warn().Msgf("Unsupported message type")
		}

		log.Info().Msgf("\n\t\tExit goroutine processingMessages(), id: %v, user: %s, name: %s %s\n", update.Message.From.ID, update.Message.From.UserName, update.Message.From.LastName, update.Message.From.FirstName)

		return
	}
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

func (bot *Bot) sendMessageToOpenaiTopic(ctx context.Context, msg *tgbotapi.Message) {

	start := time.Now()
	defer func() {
		observeResponseLatencySecSummary(time.Since(start), msg.From.UserName)
	}()

	msgText := strings.TrimPrefix(msg.Text, "/askbot ") //TODO move it to handler

	//TODO move to func
	words := strings.Split(msgText, " ")
	var filteredWords []string
	for _, word := range words {
		if !strings.HasPrefix(word, "@") {
			filteredWords = append(filteredWords, word)
		}
	}
	msgText = strings.Join(filteredWords, " ")

	openaiRequest := dto.OpenaiMsg{
		Source: dtoSource,
		ChatId: strconv.Itoa(int(msg.Chat.ID)),
		UserId: strconv.Itoa(int(msg.From.ID)),
		MsgId:  strconv.Itoa(int(msg.MessageID)),
		Msg:    msgText,
	}
	bot.broker.Publish(ctx, brokerChannelPub, openaiRequest)

}

func (bot *Bot) sendReplyMessage(inputMsg *tgbotapi.Message, replyText string) error {
	replyMsg := tgbotapi.NewMessage(inputMsg.Chat.ID, "")
	replyMsg.ReplyToMessageID = inputMsg.MessageID
	replyMsg.Text = replyText

	_, err := bot.tgbot.Send(replyMsg)
	if err != nil {
		return fmt.Errorf("tgbot.sendReplyMessage: %w", err)
	}
	return nil
}

func (bot *Bot) Shutdown(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("exit tgBot: %w", ctx.Err())
		default:
			close(bot.shutdownChannel)
			bot.cancel()
			bot.tgbot.StopReceivingUpdates()
			return nil
		}
	}
}
