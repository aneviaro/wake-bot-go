package translation

import (
	"wake-bot/usecase/translation/eng"
	"wake-bot/usecase/translation/ru"
)

// A list of constants that can be used by the translator to find a message in a particular language.
const (
	NotValidTimeFormat = iota
	ClarificationQuestion
	BestTimeToGoToSleep
	BestTimeToWakeUp
	Greetings
	NotCorrectCommand
	Usage
	GotIt
	TimeFormatQuestion
	AMPMTimeFormat
	MilitaryTimeFormat
	WakeUp
	GoToSleep
	ExpiredCallback
)

// Get finds a message by id in a langCode package.
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
