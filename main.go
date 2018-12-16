package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"webcapstaffbot/respondent"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_API_KEY"))
	if err != nil {
		log.Panic(err)
	}

	if os.Getenv("BOT_DEBUG") == "true" {
		bot.Debug = true
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		r, err := respondent.New(bot, update.Message)
		if err != nil {
			log.Println(err)
			return
		}
		r.Reply()
		//new(response.ResponseAPI).NewResponse(bot, update.Message).GenerateResponse()
	}
}
