package telegramBot

import (
	"context"
	"fmt"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/krassor/skygrow/backend-serverHttp/internal/models/dto"
	"github.com/krassor/skygrow/backend-serverHttp/internal/models/entities"
	subscriber "github.com/krassor/skygrow/backend-serverHttp/internal/services/subscriberServices"
	"github.com/rs/zerolog/log"
)

type Bot struct {
	tgbot           *tgbotapi.BotAPI
	subscriber      subscriber.SubscriberRepoService
	shutdownChannel chan struct{}
}

func NewBot(subscriber subscriber.SubscriberRepoService) *Bot {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGBOT_APITOKEN"))
	if err != nil {
		log.Error().Msgf("Error auth telegram bot: %s", err)
	}
	//TODO: add to env BOTDEBUG
	bot.Debug = false

	log.Info().Msgf("Authorized on account %s", bot.Self.UserName)

	return &Bot{
		tgbot:           bot,
		subscriber:      subscriber,
		shutdownChannel: make(chan struct{}),
	}
}

func (bot *Bot) Update(ctx context.Context, updateTimeout int) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = updateTimeout

	_, err := bot.tgbot.MakeRequest("deleteWebhook", tgbotapi.Params{"drop_pending_updates": "false"})
	if err != nil {
		log.Error().Msgf("bot.Update() error: cannot delete WebHook: %v", err)
	}

	updates := bot.tgbot.GetUpdatesChan(updateConfig)

	//TODO: make goroutine with check update channel close
	for update := range updates {

		log.Info().Msgf("Input message: %v\n", update.Message)

		if update.Message == nil && update.CallbackQuery == nil { // ignore any non-Message updates
			log.Warn().Msgf("tgbot warn: Not message: %v", update.Message)
			continue
		}

		go bot.processingMessages(update)

	}
}

func (bot *Bot) processingMessages(update tgbotapi.Update) {

	log.Info().Msgf("\n\t\tEnter goroutine processingMessages(), id: %v, user: %s, name: %s %s", update.Message.From.ID, update.Message.From.UserName, update.Message.From.LastName, update.Message.From.FirstName)
	ctx, ctxCancel := context.WithCancel(context.Background())

	select {
	case <-bot.shutdownChannel:
		ctxCancel()
	default:

		if update.Message == nil && update.CallbackQuery != nil {
			err := bot.callbackQueryHandle(ctx, update.CallbackQuery)
			if err != nil {
				log.Error().Msgf("Error tgbot handle message: %s", err)
			}
			log.Info().Msgf("CallbackQuery from user: %s, data: %s", update.CallbackQuery.From, update.CallbackQuery.Data)
		} else

		//Check if message is a command
		if update.Message.IsCommand() {
			log.Info().Msgf("tgbot.update receive command from %s: %s, text: %s", update.Message.From, update.Message.Command(), update.Message.Text)

			if err := bot.commandHandle(update.Message); err != nil {
				log.Error().Msgf("Error tgbot.update: %v", err)
			}
		} else

		// Проверяем, если сообщение адресовано самому боту
		if update.Message.Chat.IsPrivate() {
			bot.privateHandler(update.Message)
		} else

		// если сообщение адресовано каналу, в котором находится бот
		if (update.Message.Chat.IsChannel() || update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup()) && bot.checkBotMention(update.Message) {
			bot.channelHandler(update.Message)
		} else

		// Проверяем, если сообщение является ответом на сообщение бота
		if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.ID == bot.tgbot.Self.ID {
			bot.replyHandler(update.Message)
		} else {
			log.Warn().Msgf("Unsupported message type")
		}

		log.Info().Msgf("\n\t\tExit goroutine processingMessages(), id: %v, user: %s, name: %s %s\n", update.Message.From.ID, update.Message.From.UserName, update.Message.From.LastName, update.Message.From.FirstName)

		ctxCancel()
	}

}

func (bot *Bot) BookOrderNotify(ctx context.Context, bookOrder dto.ResponseBookOrderDto) error {

	msgText := fmt.Sprintf(
		"New order notify:\nBookOrderID: %v\nMentorID: %v\nFirstName: %s\nSecondName: %s\nEmail: %s\nPhone: %s\nProblem Description: %s",
		bookOrder.BookOrderID,
		bookOrder.MentorID,
		bookOrder.FirstName,
		bookOrder.SecondName,
		bookOrder.Email,
		bookOrder.Phone,
		bookOrder.ProblemDescription,
	)

	subscribers, err := bot.subscriber.GetSubscribers(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(subscribers))

	for _, s := range subscribers {
		go func(s entities.Subscriber, msgText string, wg *sync.WaitGroup) {
			msg := tgbotapi.NewMessage(s.ChatID, msgText)
			_, err := bot.tgbot.Send(msg)
			if err != nil {
				log.Error().Msgf("BookOrderNotify(): Error sending telegramm notify: %s", err)
			}
			wg.Done()
		}(s, msgText, &wg)
	}
	wg.Wait()
	return nil
}

func (bot *Bot) Shutdown(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("Force exit tgBot: %w", ctx.Err())
		default:
			close(bot.shutdownChannel)
			bot.tgbot.StopReceivingUpdates()
			return nil
		}
	}
}
