package telegramBot

import (
	"context"
	"fmt"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"github.com/krassor/skygrow/backend-serverHttp/internal/models/entities"
	services "github.com/krassor/skygrow/backend-serverHttp/internal/services/devices"
	subscriber "github.com/krassor/skygrow/backend-serverHttp/internal/services/subscriberServices"
)

type Bot struct {
	tgbot      *tgbotapi.BotAPI
	service    services.DevicesRepoService
	subscriber subscriber.SubscriberRepoService
}

func NewBot(service services.DevicesRepoService, subscriber subscriber.SubscriberRepoService) *Bot {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGBOT_APITOKEN"))
	if err != nil {
		log.Error().Msgf("Error auth telegram bot: %s", err)
	}
	//TODO: add to env BOTDEBUG
	bot.Debug = false

	log.Info().Msgf("Authorized on account %s", bot.Self.UserName)

	return &Bot{
		tgbot:      bot,
		service:    service,
		subscriber: subscriber,
	}
}

func (bot *Bot) Update(ctx context.Context, updateTimeout int) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = updateTimeout

	updates := bot.tgbot.GetUpdatesChan(updateConfig)

	//TODO: make goroutine with check update channel close
	for update := range updates {

		if update.Message == nil && update.CallbackQuery == nil { // ignore any non-Message updates
			log.Warn().Msgf("tgbot warn: Not message: %s", update.Message)
			continue
		}

		if update.Message == nil && update.CallbackQuery != nil {
			err := bot.callbackQueryHandle(ctx, update.CallbackQuery)
			if err != nil {
				log.Error().Msgf("Error tgbot handle message: %s", err)
			}
			log.Info().Msgf("CallbackQuery from user: %s, data: %s", update.CallbackQuery.From, update.CallbackQuery.Data)
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			log.Warn().Msgf("tgbot warn: Not command: %s", update.Message)
			// msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This is not command")
			// msg.ReplyToMessageID = update.Message.MessageID
			// if _, err := bot.tgbot.Send(msg); err != nil {
			// 	log.Error().Msgf("Error tgbot send message: %s", err)
			// }
			continue
		}

		log.Info().Msgf("tgbot receive command: %s", update.Message.Command())

		if err := bot.commandHandle(update.Message); err != nil {
			log.Error().Msgf("Error tgbot handle message: %s", err)
		}

	}
	log.Info().Msgf("exit telegram bot routine")
}

func (bot *Bot) commandHandle(msg *tgbotapi.Message) error {

	replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "")
	replyMsg.ReplyToMessageID = msg.MessageID

	// Extract the command from the Message.

	switch msg.Command() {
	case "help":
		replyMsg.Text = "I understand /list and /subscribe"
	case "start":
		replyMsg.Text = fmt.Sprintf("Hello, %s! I'm stand device monitor.\nEnter /list command and select device.\nEnter /subscribe to subscribe on devices status changes", msg.Chat.UserName)
	case "list":
		err := bot.list(&replyMsg)
		if err != nil {
			return err
		}
	case "subscribe":
		err := bot.subscribe(&replyMsg)
		if err != nil {
			return err
		}
	default:
		replyMsg.Text = "I don't know this command"
	}

	_, err := bot.tgbot.Send(replyMsg)
	if err != nil {
		return err
	}

	return nil
}

func (bot *Bot) DeviceStatusNotify(ctx context.Context, device entities.Devices, status bool) error {
	var statusString string
	if status {
		statusString = "ONLINE"
	} else {
		statusString = "OFFLINE"
	}

	msgText := fmt.Sprintf(
		"Device %s %s is %s",
		device.DeviceVendor,
		device.DeviceName,
		statusString,
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
				log.Error().Msgf("DeviceStatusNotify(): Error sending telegramm notify: %s", err)
			}
			wg.Done()
		}(s, msgText, &wg)
	}
	wg.Wait()
	return nil
}

// func (bot *Bot) Shutdown(ctx context.Context) error {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return fmt.Errorf("error shutdown telegram bot: %s", ctx.Err())
// 		default:
// 			bot.tgbot.StopReceivingUpdates()
// 			return nil
// 		}
// 	}
// }
