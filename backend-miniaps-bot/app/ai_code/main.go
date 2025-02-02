package main

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
	"log"
)

const (
	botToken = "5995322288:AAE9eClgD8cQyasmv5aywzVSaZ5yu-zcjcU"
	dbConn   = "user=postgres dbname=postgres sslmode=disable"
)

type Order struct {
	ID         string
	FirstName  string
	LastName   string
	MiddleName string
	Phone      string
	UserID     int64
}

func main() {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go func(update tgbotapi.Update) {
			handleUpdate(bot, update, db)
		}(update)
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *sql.DB) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Text {
	case "/start":
		msg.Text = "Привет! Нажми на кнопку, чтобы открыть форму для заказа."
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Сделать заказ", "https://cc.vtb.ru"),
			),
		)
		msg.ReplyMarkup = keyboard
	case "/status":
		var orders []Order
		_, err := db.Exec("SELECT * FROM orders WHERE user_id=$1", update.Message.From.ID)
		if err != nil {
			msg.Text = "Ошибка получения статуса заказа."
			bot.Send(msg)
			return err
		}
		if len(orders) == 0 {
			msg.Text = "У вас нет активных заказов."
		} else {
			msg.Text = "Ваши заказы:\n"
			for _, order := range orders {
				msg.Text += fmt.Sprintf("ID: %s, Имя: %s, Фамилия: %s\n", order.ID, order.FirstName, order.LastName)
			}
		}
	default:
		msg.Text = "Неизвестная команда."
	}

	_, err := bot.Send(msg)
	return err
}

func saveOrder(db *sql.DB, order Order) error {
	_, err := db.Exec(`
        INSERT INTO orders (id, first_name, last_name, middle_name, phone, user_id)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, order.ID, order.FirstName, order.LastName, order.MiddleName, order.Phone, order.UserID)
	return err
}
