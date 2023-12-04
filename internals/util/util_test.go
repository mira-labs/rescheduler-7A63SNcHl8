// File: ./internals/util/util_test.go

package util

import (
	"reflect"
	"rescheduler/internals/models"
	"rescheduler/internals/timestamp"
	"testing"
	"time"
)

func TestConvertJSONToEvent(t *testing.T) {
	// Example JSON string representing a questionnaire completed event
	jsonString := `{
		"id": "random",
		"user_id": "8a4378cd-27b9-4a36-afee-829b42eeb1b5",
		"study_id": "Study5",
		"questionnaire_id": "24b6f062-df29-4e6a-abb4-403e01671e4a",
		"completed_at": "2023-12-04 02:11:00",
		"remaining_completions": 2
	}`

	// Call the ConvertJSONToEvent function
	event, err := ConvertJSONToEvent(jsonString)

	// Check for errors
	if err != nil {
		t.Fatalf("Error converting JSON to event: %v", err)
	}

	// Validate the converted event fields
	expectedEvent := &models.QuestionnaireCompletedEvent{
		ID:                   "random",
		UserID:               "8a4378cd-27b9-4a36-afee-829b42eeb1b5",
		StudyID:              "Study5",
		QuestionnaireID:      "24b6f062-df29-4e6a-abb4-403e01671e4a",
		CompletedAt:          timestamp.TimeStamp{Time: time.Date(2023, 12, 4, 2, 11, 0, 0, time.UTC)},
		RemainingCompletions: 2,
	}

	if !reflect.DeepEqual(event, expectedEvent) {
		t.Fatalf("Converted event does not match expected event:\nGot: %+v\nExpected: %+v", event, expectedEvent)
	}
}
