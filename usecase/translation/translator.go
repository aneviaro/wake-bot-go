package translation

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
	TimezoneNotOk
	TimezoneOk
	Timezone
	SendTimezoneManually
	SendUTCOffset
	Location
	ChooseUTCOffset
	FallAsleepTimeout
)

// Get finds a message by id in a langCode package.
func Get(msgID int, langCode string) string {
	switch langCode {
	case "eng":
		return engMessages[msgID]
	case "ru":
		return ruMessages[msgID]
	default:
		return engMessages[msgID]
	}
}
