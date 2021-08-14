package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type SenderMaker interface {
	Sender
	Maker
}

type Sender interface {
	SendMessage(chatID int64, message string, replyTo int, parseMode string, keyboard *tgbotapi.
	InlineKeyboardMarkup) error
	AnswerOnCallback(id string) error
}

type Maker interface {
	MakeClarificationButtons(text1, data1, text2, data2 string) tgbotapi.InlineKeyboardMarkup
	MakeOneButton(text1, data1 string) tgbotapi.InlineKeyboardMarkup
}
