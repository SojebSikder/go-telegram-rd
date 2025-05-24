package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You said: "+update.Message.Text)

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg.Text = "Hello! I'm a bot. I can help you with your questions."
			case "help":
				msg.Text = "I can help you with your questions. Just ask me anything!"
			case "about":
				msg.Text = "Resource downloader bot"
			case "contact":
				msg.Text = "You can contact me at @sojebsikder"
			case "d":
				// get the value after /d
				args := update.Message.CommandArguments()

				if args == "" {
					msg.Text = "Usage: /say something"
				} else {
					DownloadFile(args)
					msg.Text = "File downloaded successfully"
				}
			default:
				msg.Text = "I don't know that command"
			}

			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You said: "+update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}

}
