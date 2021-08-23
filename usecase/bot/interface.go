package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// SenderMaker sets behavior of a bot methods, includes MessageSender and KeyboardMaker interfaces.
type SenderMaker interface {
	MessageSender
	KeyboardMaker
}

// MessageSender sets behaviour of a bot methods for sending messages to the bot.
type MessageSender interface {
	SendMessage(chatID int64, message string, opts ...Option) error
	AnswerOnCallback(id, label string)
	SendClarificationMessage(chatID int64, replyTo int, languageCode string) error
	SendTimeFormatMessage(chatID int64, replyTo int, languageCode string) error
	SendNotValidTimeFormatMessage(chatID int64, replyTo int, languageCode, timeFormat string) error
}

// KeyboardMaker sets behaviour of the tg keyboard manager.
type KeyboardMaker interface {
	MakeInlineKeyboard(btns ...Button) tgbotapi.InlineKeyboardMarkup
}
