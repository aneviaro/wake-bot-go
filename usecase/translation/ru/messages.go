package ru

const (
	NotValidTimeFormat     = "Неправильное время, пожалуйста, попробуйте *22:22*."
	ClarificationQuestion  = "Вы хотите лечь спать или проснутся в это время?"
	BestTimeToGoToSleep    = "Лучшее время, чтобы уснуть:\n"
	BestTimeToWakeUp       = "Лучшее время, чтобы проснуться:\n"
	Greetings              = "Привет, давайте дружить! Пожалуйста, " +
		"введите время в формате: *22:22*. Я спрошу, хотите ли вы уснуть или проснуться в это время, " +
		"и скажу вам результат. Приступим!"
	NotCorrectCommand  = "Я не знаю такой команды, она введена правильно?"
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

