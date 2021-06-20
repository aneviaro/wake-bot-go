package bot_service

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Sender interface {
	SendMessage(chatID int64, message string, replyTo int, parseMode string, keyboard interface{}) error
	MakeClarificationButtons(text1, data1, text2, data2 string) tgbotapi.InlineKeyboardMarkup
	MakeOneButton(text1, data1 string) tgbotapi.InlineKeyboardMarkup
	MakeCurrentTimeButtons(timeFormat string) tgbotapi.InlineKeyboardMarkup
	MakeRequestLocationButton() tgbotapi.ReplyKeyboardMarkup
}
