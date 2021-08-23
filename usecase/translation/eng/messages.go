package eng

// Get returns a message by id
func Get(id int) string {
	return [...]string{
		"Not a valid time format, please try *%s*.",
		"Is it a wake up time or go to sleep time?",
		"The best time to go to sleep is:\n",
		"The best time to wake up is:\n",

		"Greeting my dear sleepy friend. Good sleep is very important in our lives. " +
			"I want to make it easier for you to achieve great sleeping experience and have enough rest in you life. " +
			"My main functionality is to calculate sleep phases for you, " +
			"to make your life the most peaceful and relaxing.",

		"Apologies, I don't know this command, is it correct?",

		"Please, type in the time in the next format: " +
			"*%s*. You'll be asked to choose if it's a *Wake up* or *Go To Sleep* time. " +
			"Then I will provide the best time for you to go to sleep or wake up. Let's start!",

		"Got it!",
		"What time format you use?",
		"12-hour",
		"24-hour",
		"Wake Up",
		"Go to sleep",
		"Button expired",
	}[id]
}
