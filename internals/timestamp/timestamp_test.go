// File: ./internals/timestamp/timestamp_test.go

package timestamp

import (
	"database/sql/driver"
	"reflect"
	"testing"
	"time"
)

func TestTimeStamp_Scan(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		output TimeStamp
		err    error
	}{
		{
			name:   "Valid time string",
			input:  "2023-12-04 02:11:00",
			output: TimeStamp{Time: time.Date(2023, 12, 4, 2, 11, 0, 0, time.UTC)},
			err:    nil,
		},
		{
			name:   "Invalid time string",
			input:  "invalid_time_string",
			output: TimeStamp{},
			err:    driver.ErrSkip, // Assuming ErrSkip is returned on scan failure
		},
		{
			name:   "Nil input",
			input:  nil,
			output: TimeStamp{},
			err:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts TimeStamp
			err := ts.Scan(tt.input)

			if err != tt.err {
				t.Fatalf("Unexpected error.\nGot: %v\nExpected: %v", err, tt.err)
			}

			if !reflect.DeepEqual(ts, tt.output) {
				t.Fatalf("Unexpected TimeStamp value.\nGot: %+v\nExpected: %+v", ts, tt.output)
			}
		})
	}
}
