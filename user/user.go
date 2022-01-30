package user

import (
	"time"

	"cloud.google.com/go/datastore"
)

// TimeFormat is a string describing the possible time format for user.
type TimeFormat string

const (
	HourClock12 TimeFormat = "03:04PM"
	HourClock24 TimeFormat = "15:04"
)

// TZDefinedBy is a string describing the possible time format for user.
type TZDefinedBy string

const (
	Timezone TZDefinedBy = "timezone"
	Offset   TZDefinedBy = "offset"
)

// User is a datastore entity that represents the user.
type User struct {
	K           *datastore.Key `datastore:"__key__"`
	ChatID      int64          `json:"chat_id" datastore:"chat_id"`
	TimeFormat  TimeFormat     `json:"time_format" datastore:"time_format"`
	TimeZone    string         `json:"time_zone" datastore:"time_zone"`
	WakeUpTime  time.Time      `json:"wake_up_time" datastore:"wake_up_time"`
	UTCOffset   string         `json:"utc_offset" datastore:"utc_offset"`
	TZDefinedBy TZDefinedBy    `json:"tz_by" datastore:"tz_by"`
}

type TimeWithFormat struct {
	Format TimeFormat
	Time   time.Time
}

func (f *TimeWithFormat) String() string {
	return f.Time.Format(string(f.Format))
}
