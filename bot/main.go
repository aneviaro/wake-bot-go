package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	updatehandler "wake-bot/bot/update-handler"
	botservice "wake-bot/usecase/bot-service"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panicf("Unable to start tgbot, %v", err)
	}

	botService := botservice.MakeBotService(bot)

	handler := updatehandler.MakeUpdateHandler(botService)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		err := handler.HandleUpdate(&update)
		if err != nil {
			log.Printf("Unable to handle update, err: %v.", err)
		}
	}
}
