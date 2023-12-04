// Package database provides functionality to interact with the database for the rescheduler application.
package store

import (
	"database/sql"
	"fmt"

	"rescheduler/internals/models"
)

// QuestionnaireStoreInterface defines the methods expected for questionnaire-related database operations.
type QuestionnaireStoreInterface interface {
	FindQuestionnaireByIDAndStudyID(questionnaireID, studyID string) (*models.Questionnaire, error)
	FindQuestionnaireByID(questionnaireID string) (*models.Questionnaire, error)
}

// QuestionnaireStore implements QuestionnaireStoreInterface and is responsible for handling questionnaire-related database operations.
type QuestionnaireStore struct {
	db *sql.DB
}

// NewQuestionnaireStore creates a new QuestionnaireStore instance with the given SQL database connection.
func NewQuestionnaireStore(db *sql.DB) *QuestionnaireStore {
	return &QuestionnaireStore{db: db}
}

// FindQuestionnaireByID retrieves a questionnaire by its ID and Study ID.
// It returns a Questionnaire instance if found, or nil if no questionnaire is found.
// An error is returned if there is an issue with the database query.
// The use of studyID seemed like a contraint in the task as these details come from the event data. This shouldn't be the case as the questionnaireID should be unique therefore the key.
func (qs *QuestionnaireStore) FindQuestionnaireByIDAndStudyID(questionnaireID, studyID string) (*models.Questionnaire, error) {
	query := "SELECT * FROM questionnaires WHERE id = ? AND study_id = ?"
	row := qs.db.QueryRow(query, questionnaireID, studyID)

	var questionnaire models.Questionnaire
	err := row.Scan(
		&questionnaire.ID,
		&questionnaire.StudyID,
		&questionnaire.Name,
		&questionnaire.Questions,
		&questionnaire.MaxAttempts,
		&questionnaire.HoursBetweenAttempts,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("questionnaire not found with ID: %s and Study ID: %s", questionnaireID, studyID)
	} else if err != nil {

		return nil, err
	}

	return &questionnaire, nil
}

// FindQuestionnaireByID retrieves a questionnaire by its ID.
// It returns a Questionnaire instance if found, or nil if no questionnaire is found.
// An error is returned if there is an issue with the database query.
// This version skips studyID as it doesn't identify the questionnaire
func (qs *QuestionnaireStore) FindQuestionnaireByID(questionnaireID string) (*models.Questionnaire, error) {
	query := "SELECT * FROM questionnaires WHERE id = ?"
	row := qs.db.QueryRow(query, questionnaireID)

	var questionnaire models.Questionnaire
	err := row.Scan(
		&questionnaire.ID,
		&questionnaire.StudyID,
		&questionnaire.Name,
		&questionnaire.Questions,
		&questionnaire.MaxAttempts,
		&questionnaire.HoursBetweenAttempts,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("questionnaire not found with ID: %s", questionnaireID)
	} else if err != nil {
		return nil, fmt.Errorf("Something else is wrong %s", err)
	}

	return &questionnaire, nil
}
