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
	TimezoneNotOk:      "I can't determine your timezone.",
	TimezoneOk:         "Thanks! Your timezone is %s.\nIn case you want to change your timezone, just send me a new location via telegram attachment or in the \"UTC +5\" format.",
	Timezone: "Please, click on the button below to send your location so I can determine your timezone. " +
		"I don't store your personal information, so nothing to worry about!",
	SendTimezoneManually: "If you don't want to share your location with me, click on \"Send UTC Offset\" button below.\n" +
		"I will send you a list of available timezones.",
	SendUTCOffset:   "Send UTC Offset",
	Location:        "Current Location",
	ChooseUTCOffset: "Please, choose you UTC offset on a keyboard below.",
	FallAsleepTimeout: "The normal time it takes most people to fall asleep at night is between 10 and 20 minutes." +
		" I will calculate the best time for you to wake up with regards to this timings.",
}
