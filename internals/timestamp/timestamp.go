// Package timestamp provides a custom time type, TimeStamp, that embeds time.Time.
// It includes a Scan method for implementing the sql.Scanner interface, allowing
// seamless integration with database scanning operations.
//
// TimeStamp enhances the standard time.Time functionality and is particularly
// useful when working with databases that require custom handling of time values.
package timestamp

import (
	"encoding/json"
	"time"
)

// TimeStamp is a custom time type that embeds time.Time
// and provides a Scan method for database scanning.
type TimeStamp struct {
	time.Time
}

var layout string = "2006-01-02 15:04:05"

// Scan implements the sql.Scanner interface for TimeStamp.
// It parses the input value into a time.Time and sets it
// as the underlying time.Time field of TimeStamp.s
func (ct *TimeStamp) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	parsedTime, err := time.Parse(layout, string(value.([]byte)))
	if err != nil {
		return err
	}
	ct.Time = parsedTime
	return nil
}

func (ct *TimeStamp) UnmarshalJSON(b []byte) error {
	var timeStr string
	err := json.Unmarshal(b, &timeStr)
	if err != nil {
		return err
	}

	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		return err
	}

	ct.Time = parsedTime
	return nil
}
