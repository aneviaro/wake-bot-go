package main

import (
	"log"
	"os"
	updatehandler "wake-bot/bot/update-handler"
	"wake-bot/storage"
	botservice "wake-bot/usecase/bot-service"
	user_service "wake-bot/usecase/user-service"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panicf("Unable to start tgbot, %v", err)
	}

	store, err := storage.NewBoltStore()
	if err != nil {
		log.Panicf("Unable to start tgbot, %v", err)
	}

	botService := botservice.NewBotService(bot)
	userService := user_service.NewUserService(store)

	handler := updatehandler.MakeUpdateHandler(botService, userService)

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
