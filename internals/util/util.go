// Package util provides utility functions for the application.
package util

import (
	"encoding/json"
	"rescheduler/internals/models"
)

// ConvertJSONToEvent converts a JSON string to a models.QuestionnaireCompletedEvent.
// It uses the encoding/json package to unmarshal the JSON string into
// a models.QuestionnaireCompletedEvent struct. If successful, it returns
// a pointer to the event; otherwise, it returns an error.
func ConvertJSONToEvent(jsonString string) (*models.QuestionnaireCompletedEvent, error) {
	var event models.QuestionnaireCompletedEvent
	err := json.Unmarshal([]byte(jsonString), &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
