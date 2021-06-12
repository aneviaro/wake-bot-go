package user

import "time"

type TimeFormat string

const (
	HourClock12 TimeFormat = "03:04PM"
	HourClock24 TimeFormat = "15:04"
)

type User struct {
	ChatID     int64          `json:"chat_id"`
	TimeFormat TimeFormat     `json:"time_format"`
	TimeZone   *time.Location `json:"time_zone"`
	WakeUpTime time.Time      `json:"wake_up_time"`
}
