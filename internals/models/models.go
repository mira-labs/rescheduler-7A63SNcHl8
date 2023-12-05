/*
Package models provides the data models used in the Rescheduler application.

The models package defines the structures that represent key entities in the Rescheduler, such as participants, questionnaires, scheduled questionnaires, and questionnaire results. As I uunderstand the task, these models are designed to encapsulate the data associated with various aspects of a study, allowing for organized data representation and manipulation.

Structures:
- Participant: Represents a participant in a study with unique identification and a name.
- Questionnaire: Holds information about different questionnaires, including their configurations, maximum attempts, and scheduling parameters.
- ScheduledQuestionnaire: Represents a specific request for a participant to fill in a questionnaire at a scheduled time.
- QuestionnaireResult: Stores the results of a participant completing a questionnaire, including answers and completion timestamp.

These models in this package serve as the foundation for database interactions, providing a structured representation of the entities within the Rescheduler task.
The structure is purely based on the provided SQL script.
*/
package models

import (
	"database/sql"
	"rescheduler/internals/timestamp"
)

// Participant represents a participant in the study
type Participant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Questionnaire represents a questionnaire that participants can fill out
type Questionnaire struct {
	ID                   string        `json:"id"`
	StudyID              string        `json:"study_id"`
	Name                 string        `json:"name"`
	Questions            string        `json:"questions"`
	MaxAttempts          sql.NullInt64 `json:"max_attempts"`
	HoursBetweenAttempts int           `json:"hours_between_attempts"`
}

type ScheduledQuestionnaireStatus string

const (
	ScheduledQuestionnairePending   = "pending"
	ScheduledQuestionnaireCompleted = "completed"
)

// ScheduledQuestionnaire represents a scheduled questionnaire for a specific participant
type ScheduledQuestionnaire struct {
	ID              string                       `json:"id"`
	QuestionnaireID string                       `json:"questionnaire_id"`
	ParticipantID   string                       `json:"participant_id"`
	ScheduledAt     timestamp.TimeStamp          `json:"scheduled_at"`
	Status          ScheduledQuestionnaireStatus `json:"status"`
}

// QuestionnaireResult represents the results of a participant filling out a questionnaire
type QuestionnaireResult struct {
	ID                      string              `json:"id"`
	Answers                 string              `json:"answers"`
	QuestionnaireID         string              `json:"questionnaire_id"`
	ParticipantID           string              `json:"participant_id"`
	QuestionnaireScheduleID string              `json:"questionnaire_schedule_id"`
	CompletedAt             timestamp.TimeStamp `json:"completed_at"`
}

// QuestionnaireCompletedEvent model
// This model assumes the CompletedAt value is actually time and the whole struct will be subject to unmarshalling assuming the event comes as a json request.
// It doesn't take into consideration a Mapping Template for the Lambda integration which could ttansform each incoming property before passing it to the lambda handler.
type QuestionnaireCompletedEvent struct {
	ID                   string              `json:"id"`
	UserID               string              `json:"user_id"`
	StudyID              string              `json:"study_id"`
	QuestionnaireID      string              `json:"questionnaire_id"`
	CompletedAt          timestamp.TimeStamp `json:"completed_at"`
	RemainingCompletions int                 `json:"remaining_completions"`
}
