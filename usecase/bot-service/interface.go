package bot_service

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Sender interface {
	SendMessage(chatID int64, message string, replyTo int, parseMode string, keyboard *tgbotapi.
	InlineKeyboardMarkup) error
	MakeClarificationButtons() tgbotapi.InlineKeyboardMarkup
}
