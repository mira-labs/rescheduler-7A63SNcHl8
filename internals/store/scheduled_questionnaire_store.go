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
	SetScheduledQuestionnaireCompleted(scheduledQuestionnaire *models.ScheduledQuestionnaire) error
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
	query := "SELECT * FROM scheduled_questionnaires WHERE questionnaire_id = ? AND participant_id = ? AND status = 'pending' AND study_id = ?"
	row := scheduleStore.db.QueryRow(query, questionnaireID, userID, studyID)

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
	query := "SELECT * FROM scheduled_questionnaires WHERE questionnaire_id = ? AND participant_id = ? AND status = 'pending'"
	row := scheduleStore.db.QueryRow(query, questionnaireID, userID)
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

// SetScheduledQuestionnaireCompleted sets the status of a scheduled questionnaire to 'completed' in the database.
func (scheduleStore *ScheduledQuestionnaireStore) SetScheduledQuestionnaireCompleted(scheduledQuestionnaire *models.ScheduledQuestionnaire) error {
	query := "UPDATE scheduled_questionnaires SET status = 'completed' WHERE id = ?"
	_, err := scheduleStore.db.Exec(query, scheduledQuestionnaire.ID)
	return err
}

func (scheduleStore *ScheduledQuestionnaireStore) CreateScheduledQuestionnaire(scheduledQuestionnaire *models.ScheduledQuestionnaire) error {
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
