package user

import (
	"time"

	"cloud.google.com/go/datastore"
)

// TimeFormat is a string describing the possibnle time format for user.
type TimeFormat string

const (
	HourClock12 TimeFormat = "03:04PM"
	HourClock24 TimeFormat = "15:04"
)

// User is a datastore entity that represents the user.
type User struct {
	K          *datastore.Key `datastore:"__key__"`
	ChatID     int64          `json:"chat_id" datastore:"chat_id"`
	TimeFormat TimeFormat     `json:"time_format" datastore:"time_format"`
	TimeZone   *time.Location `json:"time_zone" datastore:"time_zone"`
	WakeUpTime time.Time      `json:"wake_up_time" datastore:"wake_up_time"`
}
