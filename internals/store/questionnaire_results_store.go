// Package store provides functionality to interact with the database for the rescheduler application.
package store

import (
	"database/sql"
	"rescheduler/internals/models"
	"time"
)

// QuestionnaireResultStoreInterface defines the methods expected for questionnaire result-related database operations.
type QuestionnaireResultStoreInterface interface {
	Create(result *models.QuestionnaireResult) error
}

// QuestionnaireResultStore implements QuestionnaireResultStoreInterface and is responsible for handling questionnaire result-related database operations.
type QuestionnaireResultStore struct {
	db *sql.DB
}

// NewQuestionnaireResultStore creates a new QuestionnaireResultStore instance with the given SQL database connection.
func NewQuestionnaireResultStore(db *sql.DB) *QuestionnaireResultStore {
	return &QuestionnaireResultStore{db: db}
}

// Create inserts a new questionnaire result record into the database.
// It takes a pointer to a QuestionnaireResult struct and inserts its values
// into the "questionnaire_results" table. The function returns an error if
// the database operation encounters any issues.
//
// Parameters:
//   - result: A pointer to a QuestionnaireResult struct containing the data to be inserted.
//
// Returns:
//   - error: An error indicating the success or failure of the database operation.
//
// Database Table Schema:
//   - Table Name: questionnaire_results
//   - Columns:
//   - id (string): Unique identifier for the questionnaire result.
//   - answers (string): JSON-encoded answers provided by the participant.
//   - questionnaire_id (string): Identifier of the associated questionnaire.
//   - participant_id (string): Identifier of the participant who completed the questionnaire.
//   - questionnaire_schedule_id (string): Identifier of the associated scheduled questionnaire.
//   - completed_at (string): Timestamp indicating when the questionnaire was completed in RFC3339 format.
func (qrs *QuestionnaireResultStore) Create(result *models.QuestionnaireResult) error {
	query := "INSERT INTO questionnaire_results (id, answers, questionnaire_id, participant_id, questionnaire_schedule_id, completed_at) VALUES (?, ?, ?, ?, ?, ?)"

	_, err := qrs.db.Exec(query, result.ID, result.Answers, result.QuestionnaireID, result.ParticipantID, result.QuestionnaireScheduleID, result.CompletedAt.Format(time.RFC3339))
	return err
}
