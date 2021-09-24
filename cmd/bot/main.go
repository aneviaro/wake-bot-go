package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"wake-bot/bot/handlers"
	"wake-bot/repository"
	botservice "wake-bot/usecase/bot"
	userservice "wake-bot/usecase/user"

	"cloud.google.com/go/datastore"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
	local := flag.Bool("local", false, "running locally")
	fromChan := flag.Bool("chan-updates", false, "getting bot updates from channel")
	fromWebhook := flag.Bool("webhook-updates", false, "getting bot updates from webhook call")

	flag.Parse()

	if *local {
		if err := godotenv.Load("local.env"); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panicf("Unable to start tgbot, err: %v", err)
	}

	dsClient, err := datastore.NewClient(context.Background(), os.Getenv("PROJECT"))
	if err != nil {
		log.Panicf("Unable to connect to datastore, err: %v", err)
	}

	defer dsClient.Close()

	store := repository.NewRepository(dsClient)

	botService := botservice.NewBotService(bot)
	userService := userservice.NewService(store)

	handler := handlers.MakeUpdateHandler(botService, userService)

	switch {
	case *fromChan:
		runChanHandler(bot, handler)
	case *fromWebhook:
		runWebhookHandler(bot, handler)
	default:
		runWebhookHandler(bot, handler)
	}
}

func runChanHandler(bot *tgbotapi.BotAPI, handler *handlers.UpdateHandler) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	_, err := bot.RemoveWebhook()
	if err != nil {
		log.Panicf("Error removing webhook: %v", err)
	}

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panicf("Error getting updates: %v.", err)
	}

	for update := range updates {
		err := handler.HandleDirectUpdate(&update)
		if err != nil {
			log.Printf("Error handling update: %v.", err)
		}
	}
}

func runWebhookHandler(bot *tgbotapi.BotAPI, handler *handlers.UpdateHandler) {
	_, err := bot.SetWebhook(tgbotapi.NewWebhook(os.Getenv("WEBHOOK_URL")))
	if err != nil {
		log.Panicf("Error setting webhook: %v.", err)
	}

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