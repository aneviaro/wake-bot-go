package update_handler

import (
	"fmt"
	"log"
	"strings"
	"time"
	botservice "wake-bot/usecase/bot-service"
	"wake-bot/usecase/translation"
	user_service "wake-bot/usecase/user-service"
	"wake-bot/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/ugjka/go-tz.v2/tz"
)

type UpdateHandler struct {
	botService  botservice.Sender
	userService user_service.Service
}

func MakeUpdateHandler(botService botservice.Sender, userService user_service.Service) *UpdateHandler {
	return &UpdateHandler{
		botService:  botService,
		userService: userService,
	}
}

func (u *UpdateHandler) HandleUpdate(update *tgbotapi.Update) error {
	if update.Message != nil {
		cmd := update.Message.Command()

		if cmd != "" {
			log.Printf("Handling incoming command: %s.", cmd)
			return u.handleCommand(update)
		}

		if update.Message.Location != nil {
			return u.handleLocationTimezone(update)
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

	markup := u.botService.MakeRequestLocationButton()

	return u.botService.SendMessage(
		update.CallbackQuery.Message.Chat.ID,
		translation.Get(translation.Timezone, langCode),
		0,
		"Markdown",
		&markup,
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

func (u *UpdateHandler) handleLocationTimezone(update *tgbotapi.Update) error {
	langCode := update.Message.From.LanguageCode

	zone, err := tz.GetZone(tz.Point{
		Lon: update.Message.Location.Longitude,
		Lat: update.Message.Location.Latitude,
	})
	if err != nil {
		e := u.botService.SendMessage(
			update.Message.Chat.ID,
			translation.Get(translation.TimezoneNotOk, langCode),
			0,
			"Markdown",
			nil,
		)

		return fmt.Errorf("%s %s", err, e)
	}

	var us user.User
	us.ChatID = update.Message.Chat.ID
	us.TimeZone = zone[0]

	err = u.userService.Update(us)
	if err != nil {
		return err
	}

	return u.botService.SendMessage(
		update.Message.Chat.ID,
		fmt.Sprintf(translation.Get(translation.TimezoneOk, langCode), zone[0]),
		0,
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
