package translation

var engMessages = map[int]string{
	NotValidTimeFormat:    "Not a valid time format, please try *%s*.",
	ClarificationQuestion: "Is it a wake up time or go to sleep time?",
	BestTimeToGoToSleep:   "The best time to go to sleep is:\n",
	BestTimeToWakeUp:      "The best time to wake up is:\n",
	Greetings: "Greeting my dear sleepy friend. Good sleep is very important in our lives. " +
		"I want to make it easier for you to achieve great sleeping experience and have enough rest in you life. " +
		"My main functionality is to calculate sleep phases for you, " +
		"to make your life the most peaceful and relaxing.",
	NotCorrectCommand: "Apologies, I don't know this command, is it correct?",
	Usage: "Please, type in the time in the next format: " +
		"*%s*. You'll be asked to choose if it's a *Wake up* or *Go To Sleep* time. " +
		"Then I will provide the best time for you to go to sleep or wake up. Let's start!",
	GotIt:              "Got it!",
	TimeFormatQuestion: "What time format you use?",
	AMPMTimeFormat:     "12-hour",
	MilitaryTimeFormat: "24-hour",
	WakeUp:             "Wake Up",
	GoToSleep:          "Go to sleep",
	ExpiredCallback:    "Button expired",
}
