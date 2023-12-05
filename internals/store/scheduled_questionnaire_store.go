// Package database provides functionality to interact with the database for the rescheduler application.
package store

import (
	"database/sql"
	"fmt"
	"rescheduler/internals/models"
	"time"
)

// ScheduledQuestionnaireStoreInterface defines the methods expected for scheduled questionnaire-related database operations.
type ScheduledQuestionnaireStoreInterface interface {
	FindScheduledQuestionnaireByQuestionnaireIDAndUserIDAndStudyID(questionnaireID, userID, studyID string) (*models.ScheduledQuestionnaire, error)
	FindScheduledQuestionnaireByQuestionnaireIDAndUserID(questionnaireID string, userID string) (*models.ScheduledQuestionnaire, error)
	Update(scheduledQuestionnaire *models.ScheduledQuestionnaire) error
	Create(scheduledQuestionnaire *models.ScheduledQuestionnaire) error
}

// ScheduledQuestionnaireStore implements ScheduledQuestionnaireStoreInterface and is responsible for handling scheduled questionnaire-related database operations.
type ScheduledQuestionnaireStore struct {
	db *sql.DB
}

// NewScheduledQuestionnaireStore creates a new ScheduledQuestionnaireStore instance with the given SQL database connection.
func NewScheduledQuestionnaireStore(db *sql.DB) *ScheduledQuestionnaireStore {
	return &ScheduledQuestionnaireStore{db: db}
}

// FindScheduledQuestionnaireByQuestionnaireIDAndUserIDAndStudyID retrieves a scheduled questionnaire by QuestionnaireID, UserID, and StudyID with a specific status.
// It returns a ScheduledQuestionnaire instance if found, or nil if no scheduled questionnaire is found.
// An error is returned if there is an issue with the database query.
// This again shouldn't need the studyID as it doesn't identify the questionnaire but it comes with the event data
func (scheduleStore *ScheduledQuestionnaireStore) FindScheduledQuestionnaireByQuestionnaireIDAndUserIDAndStudyID(questionnaireID, userID, studyID string) (*models.ScheduledQuestionnaire, error) {
	query := "SELECT * FROM scheduled_questionnaires WHERE questionnaire_id = ? AND participant_id = ? AND status =?  AND study_id = ?"
	row := scheduleStore.db.QueryRow(query, questionnaireID, userID, models.ScheduledQuestionnairePending, studyID)

	var scheduledQuestionnaire models.ScheduledQuestionnaire
	err := row.Scan(
		&scheduledQuestionnaire.ID,
		&scheduledQuestionnaire.QuestionnaireID,
		&scheduledQuestionnaire.ParticipantID,
		&scheduledQuestionnaire.ScheduledAt,
		&scheduledQuestionnaire.Status,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("scheduled questionnaire not found with Questionnaire ID: %s, User ID: %s, and Study ID: %s", questionnaireID, userID, studyID)
	} else if err != nil {
		return nil, err
	}

	return &scheduledQuestionnaire, nil
}

// FindScheduledQuestionnaireByQuestionnaireIDAndUserIDAndStudyID retrieves a pending scheduled questionnaire by QuestionnaireID, UserID
// It returns a ScheduledQuestionnaire instance if found, or nil if no scheduled questionnaire is found.
// An error is returned if there is an issue with the database query.
// This version omits studyID as this doesn't identify a scheduledQuestionnaire
func (scheduleStore *ScheduledQuestionnaireStore) FindScheduledQuestionnaireByQuestionnaireIDAndUserID(questionnaireID string, userID string) (*models.ScheduledQuestionnaire, error) {
	query := "SELECT * FROM scheduled_questionnaires WHERE questionnaire_id = ? AND participant_id = ? AND status = ?"
	row := scheduleStore.db.QueryRow(query, questionnaireID, userID, models.ScheduledQuestionnairePending)
	var scheduledQuestionnaire models.ScheduledQuestionnaire
	err := row.Scan(
		&scheduledQuestionnaire.ID,
		&scheduledQuestionnaire.QuestionnaireID,
		&scheduledQuestionnaire.ParticipantID,
		&scheduledQuestionnaire.ScheduledAt,
		&scheduledQuestionnaire.Status,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("scheduled questionnaire not found with Questionnaire ID: %s, User ID: %s", questionnaireID, userID)
	} else if err != nil {
		return nil, err
	}

	return &scheduledQuestionnaire, nil
}

func (scheduleStore *ScheduledQuestionnaireStore) Update(scheduledQuestionnaire *models.ScheduledQuestionnaire) error {
	query := "UPDATE scheduled_questionnaires SET questionnaire_id = ?, participant_id = ?, scheduled_at = ?, status = ? WHERE id = ?"
	_, err := scheduleStore.db.Exec(query, scheduledQuestionnaire.QuestionnaireID, scheduledQuestionnaire.ParticipantID, scheduledQuestionnaire.ScheduledAt, scheduledQuestionnaire.Status, scheduledQuestionnaire.ID)
	return err
}

// Create inserts a new scheduled questionnaire record into the database.
// It takes a pointer to a ScheduledQuestionnaire struct and inserts its values
// into the "scheduled_questionnaires" table. The function returns an error if
// the database operation encounters any issues.
//
// Parameters:
//   - scheduledQuestionnaire: A pointer to a ScheduledQuestionnaire struct
//     containing the data to be inserted.
//
// Returns:
//   - error: An error indicating the success or failure of the database operation.
//
// Database Table Schema:
//   - Table Name: scheduled_questionnaires
//   - Columns:
//   - id (string): Unique identifier for the scheduled questionnaire.
//   - questionnaire_id (string): Identifier of the associated questionnaire.
//   - participant_id (string): Identifier of the participant assigned to the questionnaire.
//   - scheduled_at (string): Scheduled time of the questionnaire in RFC3339 format.
//   - status (string): Status of the scheduled questionnaire (e.g., "pending" or "completed").
func (scheduleStore *ScheduledQuestionnaireStore) Create(scheduledQuestionnaire *models.ScheduledQuestionnaire) error {
	query := "INSERT INTO scheduled_questionnaires (id, questionnaire_id, participant_id, scheduled_at, status) VALUES (?, ?, ?, ?, ?)"

	_, err := scheduleStore.db.Exec(query,
		scheduledQuestionnaire.ID,
		scheduledQuestionnaire.QuestionnaireID,
		scheduledQuestionnaire.ParticipantID,
		scheduledQuestionnaire.ScheduledAt.Format(time.RFC3339),
		scheduledQuestionnaire.Status,
	)

	return err
}
