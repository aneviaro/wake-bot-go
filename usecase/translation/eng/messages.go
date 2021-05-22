package eng

const (
	NotValidTimeFormat     = "Not a valid time format, please try *22:22*."
	ClarificationQuestion  = "Is it wake up time or go to sleep time?"
	BestTimeToGoToSleep    = "The best time to go to sleep is:\n"
	BestTimeToWakeUp       = "The best time to wake up is:\n"
	Greetings              = "Greeting my dear sleepy friend. Let's be friends! Please, " +
		"type in the time in the next format: I.e. " +
		"*22:15*. You'll be asked to choose if it's a *Wake up* or *Go To Sleep* time. " +
		"Then I will provide the best time for you to go to sleep or wake up. Let's start!"
	NotCorrectCommand  = "Apologies, I don't know this command, is it correct?"
)

func Get(id int) string {
	switch id {
	case 0:
		return NotValidTimeFormat
	case 1:
		return ClarificationQuestion
	case 2:
		return BestTimeToGoToSleep
	case 3:
		return BestTimeToWakeUp
	case 4:
		return Greetings
	case 5:
		return NotCorrectCommand
	default:
		return ""
	}
}
