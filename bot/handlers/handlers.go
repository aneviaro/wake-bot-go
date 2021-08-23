package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	botservice "wake-bot/usecase/bot"
	"wake-bot/usecase/callback"
	"wake-bot/usecase/translation"
	userservice "wake-bot/usecase/user"
	"wake-bot/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// UpdateHandler is handler for updates that comes from the Telegram bot.
type UpdateHandler struct {
	botService  botservice.SenderMaker
	userService userservice.IService
}

// MakeUpdateHandler creates a new UpdateHandler.
func MakeUpdateHandler(botService botservice.SenderMaker, userService userservice.IService) *UpdateHandler {
	return &UpdateHandler{
		botService:  botService,
		userService: userService,
	}
}

// parseTelegramRequest decodes telegram web hook request body
func parseTelegramRequest(r *http.Request) (*tgbotapi.Update, error) {
	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		return nil, err
	}

	return &update, nil
}

// HandleTelegramWebHook decodes and handles telegram web hook request.
func (u *UpdateHandler) HandleTelegramWebHook(_ http.ResponseWriter, r *http.Request) {
	var update, err = parseTelegramRequest(r)
	if err != nil {
		log.Printf("Met error while parsing update, err: %s", err.Error())
		return
	}

	err = u.handleUpdate(update)
	if err != nil {
		log.Printf("Unable to handle update, err: %v.", err)
	}
}

// HandleDirectUpdate handles the direct message from updates channel.
func (u *UpdateHandler) HandleDirectUpdate(update *tgbotapi.Update) error {
	return u.handleUpdate(update)
}

// handleUpdate checks the update type and handles it
func (u *UpdateHandler) handleUpdate(update *tgbotapi.Update) error {
	if update.Message != nil {
		if cmd := update.Message.Command(); cmd != "" {
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

// handleMessage handles message update type
func (u *UpdateHandler) handleMessage(update *tgbotapi.Update) error {
	var timeFormat user.TimeFormat

	chatID, langCode := update.Message.Chat.ID, update.Message.From.LanguageCode

	us, err := u.userService.GetByID(chatID)
	if err != nil || us.ChatID == 0 {
		return u.botService.SendTimeFormatMessage(chatID, 0, langCode)
	}

	timeFormat = us.TimeFormat

	_, err = time.Parse(string(timeFormat), strings.Trim(update.Message.Text, " .,"))
	if err != nil {
		return u.botService.SendNotValidTimeFormatMessage(chatID, 0, langCode, string(timeFormat))
	}

	return u.botService.SendClarificationMessage(update.Message.Chat.ID, update.Message.MessageID, langCode)
}

// handleCallback handles clarification, gotit and time format callback types.
func (u *UpdateHandler) handleCallback(update *tgbotapi.Update) error {
	var err error

	callbackData := update.CallbackQuery.Data
	callbackLabel := ""

	switch {
	case callback.IsClarification(callbackData):
		err = u.handleClarificationCallback(update)
	case callback.IsGotIt(callbackData):
		err = u.handleGotItCallback(update)
	case callback.IsTimeSelect(callbackData):
		err = u.handleTimeFormatCallback(update)
	default:
		callbackLabel = translation.Get(translation.ExpiredCallback, update.CallbackQuery.Message.From.LanguageCode)
	}

	// prevents multiple answers on same callback
	go u.botService.AnswerOnCallback(update.CallbackQuery.ID, callbackLabel)

	return err
}

// handleGotItCallback handles got it callback.
func (u *UpdateHandler) handleGotItCallback(update *tgbotapi.Update) error {
	langCode := update.CallbackQuery.From.LanguageCode
	chatID := update.CallbackQuery.Message.Chat.ID

	return u.botService.SendTimeFormatMessage(chatID, 0, langCode)
}

// handleTimeFormatCallback handles the time format callback.
func (u *UpdateHandler) handleTimeFormatCallback(update *tgbotapi.Update) error {
	langCode := update.CallbackQuery.From.LanguageCode

	var us user.User
	us.ChatID = update.CallbackQuery.Message.Chat.ID

	switch update.CallbackQuery.Data {
	case callback.AMPMTime:
		us.TimeFormat = user.HourClock12
	case callback.MilitaryTime:
		us.TimeFormat = user.HourClock24
	default:
		us.TimeFormat = user.HourClock12
	}

	err := u.userService.Update(us)
	if err != nil {
		return err
	}

	return u.botService.SendMessage(
		update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf(translation.Get(translation.Usage, langCode), time.Now().UTC().Format(string(us.TimeFormat))),
	)
}

// handleCommand handles the commands.
func (u *UpdateHandler) handleCommand(update *tgbotapi.Update) error {
	langCode := update.Message.From.LanguageCode
	if update.Message.Text == "/start" || strings.Contains(update.Message.Text, "/restart") {
		gotItButton := u.botService.MakeInlineKeyboard(botservice.NewButton(translation.Get(translation.GotIt, langCode), callback.GotIt))

		err := u.botService.SendMessage(
			update.Message.Chat.ID,
			translation.Get(translation.Greetings, langCode),
			botservice.WithKeyboard(&gotItButton),
		)

		return err
	}

	return u.botService.SendMessage(
		update.Message.Chat.ID,
		translation.Get(translation.NotCorrectCommand, langCode),
		botservice.WithReplyTo(update.Message.MessageID),
	)
}

// handleClarificationCallback handles the clarification callback.
func (u *UpdateHandler) handleClarificationCallback(update *tgbotapi.Update) error {
	clarificationAnswer := update.CallbackQuery.Data
	replyToMsg, langCode := "", update.CallbackQuery.Message.From.LanguageCode

	if update.CallbackQuery.Message.ReplyToMessage != nil {
		replyToMsg, langCode = update.CallbackQuery.Message.ReplyToMessage.Text, update.CallbackQuery.Message.ReplyToMessage.From.LanguageCode
	}

	chatID := update.CallbackQuery.Message.Chat.ID

	var timeFormat user.TimeFormat

	us, err := u.userService.GetByID(update.CallbackQuery.Message.Chat.ID)
	if err != nil {
		timeFormat = user.HourClock12
	} else {
		timeFormat = us.TimeFormat
	}

	t, err := time.Parse(string(timeFormat), strings.Trim(replyToMsg, " .,"))
	if err != nil {
		return u.botService.SendNotValidTimeFormatMessage(chatID, 0, langCode, string(timeFormat))
	}

	var times []time.Time

	var msgTest string

	switch clarificationAnswer {
	case callback.WakeUp:
		msgTest = translation.Get(translation.BestTimeToGoToSleep, langCode)
		times = makeTimesArr(&t, -1)
	case callback.GoToSleep:
		msgTest = translation.Get(translation.BestTimeToWakeUp, langCode)
		times = makeTimesArr(&t, 1)
	}

	for _, t := range times {
		msgTest += fmt.Sprintf("`%s`\n", t.Format(string(timeFormat)))
	}

	return u.botService.SendMessage(update.CallbackQuery.Message.Chat.ID, msgTest)
}

// makeTimesArr creates a list of times to be sent to the bot.
func makeTimesArr(inputTime *time.Time, coef int) []time.Time {
	var times []time.Time
	for i := 0; i < 6; i++ {
		times = append(times, inputTime.Add(time.Duration(coef*90*(i+1))*time.Minute))
	}

	return times
}
