package main

import (
	"context"
	"log"
	"net/http"
	"os"
	updatehandler "wake-bot/bot/update-handler"
	"wake-bot/storage"
	botservice "wake-bot/usecase/bot-service"
	user_service "wake-bot/usecase/user-service"

	"cloud.google.com/go/datastore"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panicf("Unable to start tgbot, err: %v", err)
	}

	ctx := context.Background()
	dsClient, err := datastore.NewClient(ctx, os.Getenv("PROJECT"))
	if err != nil {
		log.Panicf("Unable to connect to datastore, err: %v", err)
	}

	defer dsClient.Close()

	store := storage.NewDatastore(dsClient)

	botService := botservice.NewBotService(bot)
	userService := user_service.NewUserService(store)

	handler := updatehandler.MakeUpdateHandler(botService, userService)
	http.HandleFunc("/", handler.HandleTelegramWebHook)


	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Panicln(err)
	}
}