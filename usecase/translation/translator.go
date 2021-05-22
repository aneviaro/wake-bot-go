package translation

import (
	"wake-bot/usecase/translation/eng"
	"wake-bot/usecase/translation/ru"
)

const (
	NotValidTimeFormat = iota
	ClarificationQuestion
	BestTimeToGoToSleep
	BestTimeToWakeUp
	Greetings
	NotCorrectCommand
)

func Get(msgID int, langCode string) string {
	switch langCode {
	case "eng":
		return eng.Get(msgID)
	case "ru":
		return ru.Get(msgID)
	default:
		return eng.Get(msgID)
	}
}
