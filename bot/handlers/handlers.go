package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	botservice "wake-bot/usecase/bot"
	"wake-bot/usecase/translation"
	userservice "wake-bot/usecase/user"
	"wake-bot/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type UpdateHandler struct {
	botService  botservice.SenderMaker
	userService userservice.IService
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func parseTelegramRequest(r *http.Request) (*tgbotapi.Update, error) {
	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return nil, err
	}
	return &update, nil
}

func (u *UpdateHandler) HandleTelegramWebHook(_ http.ResponseWriter, r *http.Request) {
	var update, err = parseTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

	err = u.handleUpdate(update)
	if err != nil {
		log.Printf("Unable to handle update, err: %v.", err)
	}
}

func MakeUpdateHandler(botService botservice.SenderMaker, userService userservice.IService) *UpdateHandler {
	return &UpdateHandler{
		botService:  botService,
		userService: userService,
	}
}

func (u *UpdateHandler) handleUpdate(update *tgbotapi.Update) error {
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

		// prevents multiple answers on same callback
		_ = u.botService.AnswerOnCallback(update.CallbackQuery.ID)
		return u.handleCallback(update)
	}

	return nil
}

func (u *UpdateHandler) handleMessage(update *tgbotapi.Update) error {

	us, err := u.userService.GetByID(update.Message.Chat.ID)
	var timeFormat user.TimeFormat

	if err != nil || us.ChatID == 0 {
		return u.sendTimeFormatQuestion(update.Message.From.LanguageCode, update.Message.Chat.ID)
	}

	timeFormat = us.TimeFormat

	_, err = time.Parse(string(timeFormat), strings.Trim(update.Message.Text, " .,"))
	if err != nil {
		return u.botService.SendMessage(
			update.Message.Chat.ID,
			fmt.Sprintf(translation.Get(translation.NotValidTimeFormat, update.Message.From.LanguageCode),
				time.Now().UTC().Format(string(timeFormat))),
			0,
			"Markdown",
			nil,
		)
	}

	keyboard := u.botService.MakeClarificationButtons(
		translation.Get(translation.WakeUp, update.Message.From.LanguageCode),
		"clarification1",
		translation.Get(translation.GoToSleep, update.Message.From.LanguageCode),
		"clarification2",
	)

	return u.botService.SendMessage(
		update.Message.Chat.ID,
		translation.Get(translation.ClarificationQuestion, update.Message.From.LanguageCode),
		update.Message.MessageID,
		"Markdown",
		&keyboard,
	)
}

func (u *UpdateHandler) handleCallback(update *tgbotapi.Update) error {
	callbackData := update.CallbackQuery.Data
	var err error

	switch {
	case strings.Contains(callbackData, "clarification"):
		err = u.handleClarificationCallback(update)
	case strings.Contains(callbackData, "gotit"):
		err = u.handleGotItCallback(update)
	case strings.Contains(callbackData, "timeFormat"):
		err = u.handleTimeFormatCallback(update)
	}

	return err
}

func (u UpdateHandler) handleGotItCallback(update *tgbotapi.Update) error {
	langCode := update.CallbackQuery.From.LanguageCode
	chatID := update.CallbackQuery.Message.Chat.ID
	return u.sendTimeFormatQuestion(langCode, chatID)
}

func (u UpdateHandler) sendTimeFormatQuestion(langCode string, chatID int64) error {
	timeFormatButtons := u.botService.MakeClarificationButtons(
		translation.Get(translation.TimeFormat1, langCode),
		"timeFormat1",
		translation.Get(translation.TimeFormat2, langCode),
		"timeFormat2",
	)

	return u.botService.SendMessage(
		chatID,
		translation.Get(translation.TimeFormatQuestion, langCode),
		0,
		"Markdown",
		&timeFormatButtons,
	)
}

func (u UpdateHandler) handleTimeFormatCallback(update *tgbotapi.Update) error {
	langCode := update.CallbackQuery.From.LanguageCode

	var us user.User
	us.ChatID = update.CallbackQuery.Message.Chat.ID
	switch update.CallbackQuery.Data {
	case "timeFormat1":
		us.TimeFormat = user.HourClock12
	case "timeFormat2":
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
		0,
		"Markdown",
		nil,
	)
}

func (u UpdateHandler) handleCommand(update *tgbotapi.Update) error {
	langCode := update.Message.From.LanguageCode
	if update.Message.Text == "/start" || strings.Contains(update.Message.Text, "/restart") {
		gotItButton := u.botService.MakeOneButton(translation.Get(translation.GotIt, langCode), "gotit")
		err := u.botService.SendMessage(
			update.Message.Chat.ID,
			translation.Get(translation.Greetings, langCode),
			0,
			"Markdown",
			&gotItButton,
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

func (u *UpdateHandler) handleClarificationCallback(update *tgbotapi.Update) error {
	clarificationAnswer := update.CallbackQuery.Data
	replyToMsg := update.CallbackQuery.Message.ReplyToMessage.Text
	langCode := update.CallbackQuery.Message.ReplyToMessage.From.LanguageCode

	var timeFormat user.TimeFormat

	us, err := u.userService.GetByID(update.CallbackQuery.Message.Chat.ID)
	if err != nil {
		timeFormat = user.HourClock12
	} else {
		timeFormat = us.TimeFormat
	}

	t, err := time.Parse(string(timeFormat), strings.Trim(replyToMsg, " .,"))
	if err != nil {
		return u.botService.SendMessage(
			update.Message.Chat.ID,
			fmt.Sprintf(translation.Get(translation.NotValidTimeFormat, update.Message.From.LanguageCode),
				time.Now().UTC().Format(string(timeFormat))),
			0,
			"Markdown",
			nil,
		)
	}

	var times []time.Time
	var msgTest string

	switch clarificationAnswer {
	case "clarification1":
		msgTest = translation.Get(translation.BestTimeToGoToSleep, langCode)

		times = makeTimesArr(&t, -1)
	case "clarification2":
		msgTest = translation.Get(translation.BestTimeToWakeUp, langCode)

		times = makeTimesArr(&t, 1)
	}

	for _, t := range times {
		msgTest += fmt.Sprintf("`%s`\n", t.Format(string(timeFormat)))
	}

	return u.botService.SendMessage(update.CallbackQuery.Message.Chat.ID, msgTest, 0,
		"Markdown", nil)
}

func makeTimesArr(inputTime *time.Time, coef int) []time.Time {
	var times []time.Time
	for i := 0; i < 6; i++ {
		times = append(times, inputTime.Add(time.Duration(coef*90*(i+1))*time.Minute))
	}

	return times
}