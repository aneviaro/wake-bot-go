package bot_service

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type BotService struct {
	tg *tgbotapi.BotAPI
}

func NewBotService(bot *tgbotapi.BotAPI) *BotService {
	return &BotService{tg: bot}
}

func (bot *BotService) SendMessage(chatID int64, message string, replyTo int, parseMode string,
	keyboard *tgbotapi.InlineKeyboardMarkup) error {
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
