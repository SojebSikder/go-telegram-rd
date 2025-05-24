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

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello! I'm a bot. I can help you with your questions.")
				bot.Send(msg)
			case "help":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I can help you with your questions. Just ask me anything!")
				bot.Send(msg)
			case "about":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Resource downloader bot created by @sojebsikder")
				bot.Send(msg)
			case "contact":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You can contact me at @sojebsikder")
				bot.Send(msg)
			case "d":
				// get the value after /d
				args := update.Message.CommandArguments()

				if args == "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: /d <url>")
					bot.Send(msg)
				} else {
					free := IsResourceFree(args)

					if free {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This is already free. Don't waste my time.")
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					} else {
						file, filename, _ := DownloadFile(args, true)

						// to group
						privateMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "File sent to you privately")
						bot.Send(privateMsg)

						// show download progress to user
						progress := tgbotapi.NewMessage(update.Message.From.ID, "Downloading...")
						bot.Send(progress)

						// send file to user
						privateFile := tgbotapi.NewDocument(update.Message.From.ID, tgbotapi.FileBytes{Name: filename, Bytes: file})
						bot.Send(privateFile)
					}

				}
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
				bot.Send(msg)
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You said: "+update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}

}
