package bot_service

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotService struct {
	tg *tgbotapi.BotAPI
}

func NewBotService(bot *tgbotapi.BotAPI) *BotService {
	return &BotService{tg: bot}
}

func (bot *BotService) SendMessage(chatID int64, message string, replyTo int, parseMode string,
	keyboard interface{}) error {
	msg := tgbotapi.NewMessage(chatID, message)
	if replyTo != 0 {
		msg.ReplyToMessageID = replyTo
	}

	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}

	msg.ParseMode = parseMode
	_, err := bot.tg.Send(msg)
	return err
}

func (bot *BotService) MakeClarificationButtons(text1, data1, text2, data2 string) tgbotapi.InlineKeyboardMarkup {
	var buttons []tgbotapi.InlineKeyboardButton
	buttons = append(buttons,
		tgbotapi.NewInlineKeyboardButtonData(text1, data1),
		tgbotapi.NewInlineKeyboardButtonData(text2, data2),
	)

	return tgbotapi.NewInlineKeyboardMarkup(buttons)
}

func (bot *BotService) MakeOneButton(text1, data1 string) tgbotapi.InlineKeyboardMarkup {
	var buttons []tgbotapi.InlineKeyboardButton
	buttons = append(buttons,
		tgbotapi.NewInlineKeyboardButtonData(text1, data1),
	)
	return tgbotapi.NewInlineKeyboardMarkup(buttons)
}

func (bot *BotService) MakeCurrentTimeButtons(timeFormat string) tgbotapi.InlineKeyboardMarkup {
	now := time.Now()
	currentMinutes := now.Minute()
	startOfToday := time.Now().Truncate(24 * time.Hour).Add(time.Duration(currentMinutes) * time.Minute)
	buttons := make([]tgbotapi.InlineKeyboardButton, 24)
	for i := 0; i < 24; i++ {
		startOfHour := startOfToday.Add(1 * time.Hour).Format(timeFormat)
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(startOfHour, fmt.Sprintf("timezone.%d", i))
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons)
}

func (bot BotService) MakeRequestLocationButton() tgbotapi.ReplyKeyboardMarkup {
	button := tgbotapi.NewKeyboardButton("Location/Timezone")
	button.RequestLocation = true

	markup := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(button))
	markup.ResizeKeyboard = true
	markup.OneTimeKeyboard = true

	return markup
}
