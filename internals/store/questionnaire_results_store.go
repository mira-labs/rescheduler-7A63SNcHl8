// Package store provides functionality to interact with the database for the rescheduler application.
package store

import (
	"database/sql"
	"rescheduler/internals/models"
	"time"
)

// QuestionnaireResultStoreInterface defines the methods expected for questionnaire result-related database operations.
type QuestionnaireResultStoreInterface interface {
	CreateQuestionnaireResult(result *models.QuestionnaireResult) error
}

// QuestionnaireResultStore implements QuestionnaireResultStoreInterface and is responsible for handling questionnaire result-related database operations.
type QuestionnaireResultStore struct {
	db *sql.DB
}

// NewQuestionnaireResultStore creates a new QuestionnaireResultStore instance with the given SQL database connection.
func NewQuestionnaireResultStore(db *sql.DB) *QuestionnaireResultStore {
	return &QuestionnaireResultStore{db: db}
}

// CreateQuestionnaireResult creates a new questionnaire result in the database.
func (qrs *QuestionnaireResultStore) CreateQuestionnaireResult(result *models.QuestionnaireResult) error {

	query := "INSERT INTO questionnaire_results (id, answers, questionnaire_id, participant_id, questionnaire_schedule_id, completed_at) VALUES (?, ?, ?, ?, ?, ?)"

	_, err := qrs.db.Exec(query, result.ID, result.Answers, result.QuestionnaireID, result.ParticipantID, result.QuestionnaireScheduleID, result.CompletedAt.Format(time.RFC3339))
	return err
}
