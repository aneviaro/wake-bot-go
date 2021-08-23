package callback

// A set of constants that are to setup possible callbacks.
const (
	WakeUp       string = "wu_cl"
	GoToSleep    string = "gts_cl"
	MilitaryTime string = "mt_cl"
	AMPMTime     string = "ampm_cl"
	GotIt        string = "gotit_cl"
)

// IsClarification checks if the callback string is one of the clarification answers.
func IsClarification(str string) bool {
	return str == WakeUp || str == GoToSleep
}

// IsTimeSelect checks if the callback string is one of the time select answers.
func IsTimeSelect(str string) bool {
	return str == MilitaryTime || str == AMPMTime
}

// IsGotIt checks if the callback string is a gotit answer.
func IsGotIt(str string) bool {
	return str == GotIt
}
