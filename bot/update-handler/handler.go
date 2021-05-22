package update_handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"time"
	botservice "wake-bot/usecase/bot-service"
	"wake-bot/usecase/translation"
)

type UpdateHandler struct {
	botService botservice.Sender
}

func MakeUpdateHandler(botService botservice.Sender) *UpdateHandler {
	return &UpdateHandler{
		botService: botService,
	}
}

func (u *UpdateHandler) HandleUpdate(update *tgbotapi.Update) error {
	if update.Message != nil {
		cmd := update.Message.Command()

		if cmd != "" {
			log.Printf("Handling incoming command: %s.", cmd)
			return u.handleCommand(update)
		}

		log.Printf("Handling incoming message: %s.", update.Message.Text)
		return u.handleMessage(update)
	}

	if update.CallbackQuery != nil {
		log.Printf("Handling incoming callback: %s.", update.CallbackQuery.Data)
		return u.handleCallback(update)
	}

	return nil
}

func (u *UpdateHandler) handleMessage(update *tgbotapi.Update) error {

	_, err := time.Parse("15:04", strings.Trim(update.Message.Text, " .,"))
	if err != nil {
		err := u.botService.SendMessage(
			update.Message.Chat.ID,
			translation.Get(translation.NotValidTimeFormat, update.Message.From.LanguageCode),
			0,
			"Markdown",
			nil,
		)
		return err
	}

	keyboard := u.botService.MakeClarificationButtons()

	return u.botService.SendMessage(
		update.Message.Chat.ID,
		translation.Get(translation.ClarificationQuestion, update.Message.From.LanguageCode),
		update.Message.MessageID,
		"Markdown",
		&keyboard,
	)
}

func (u *UpdateHandler) handleCallback(update *tgbotapi.Update) error {
	clarificationAnswer := update.CallbackQuery.Data
	replyToMsg := update.CallbackQuery.Message.ReplyToMessage.Text
	langCode := update.CallbackQuery.Message.ReplyToMessage.From.LanguageCode

	t, err := time.Parse("15:04", strings.Trim(replyToMsg, " .,"))
	if err != nil {
		err := u.botService.SendMessage(
			update.Message.Chat.ID,
			translation.Get(translation.NotValidTimeFormat, langCode),
			0,
			"Markdown",
			nil,
		)
		return err
	}

	var times []time.Time
	var msgTest string

	if clarificationAnswer == "1" {
		msgTest = translation.Get(translation.BestTimeToGoToSleep, langCode)

		times = makeTimesArr(&t, -1)
	} else if clarificationAnswer == "2" {
		msgTest = translation.Get(translation.BestTimeToWakeUp, langCode)

		times = makeTimesArr(&t, 1)
	}

	for _, t := range times {
		msgTest += fmt.Sprintf("`%s`\n", t.Format("15:04"))
	}

	return u.botService.SendMessage(update.CallbackQuery.Message.Chat.ID, msgTest, 0,
		"Markdown", nil)
}

func (u UpdateHandler) handleCommand(update *tgbotapi.Update) error {
	langCode := update.Message.From.LanguageCode
	if update.Message.Text == "/start" || strings.Contains(update.Message.Text, "/restart") {
		err := u.botService.SendMessage(
			update.Message.Chat.ID,
			translation.Get(translation.Greetings, langCode),
			0,
			"Markdown",
			nil,
		)
		return err
	}

	return u.botService.SendMessage(
		update.Message.Chat.ID,
		translation.Get(translation.NotCorrectCommand, langCode),
		update.Message.MessageID,
		"Markdown",
		nil,
	)
}

func makeTimesArr(inputTime *time.Time, coef int) []time.Time {
	var times []time.Time
	for i := 0; i < 6; i++ {
		times = append(times, inputTime.Add(time.Duration(coef*90*(i+1))*time.Minute))
	}

	return times
}
