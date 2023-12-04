// Package database provides functionality to interact with the database for the rescheduler application.
package store

import (
	"database/sql"
	"fmt"

	"rescheduler/internals/models"
)

// ParticipantStoreInterface defines the methods expected for participant-related database operations.
type ParticipantStoreInterface interface {
	FindParticipantByID(participantID string) (*models.Participant, error)
}

// ParticipantStore implements ParticipantStoreInterface and is responsible for handling participant-related database operations.
type ParticipantStore struct {
	db *sql.DB
}

// NewParticipantStore creates a new ParticipantStore instance with the given SQL database connection.
func NewParticipantStore(db *sql.DB) *ParticipantStore {
	return &ParticipantStore{db: db}
}

// FindParticipantByID retrieves a participant by their ID.
// It returns a Participant instance if found, or nil if no participant is found.
// An error is returned if there is an issue with the database query.
// So far this function is unused but I creqted it for the sake of consistency in maintaining the store concept.
func (ps *ParticipantStore) FindParticipantByID(participantID string) (*models.Participant, error) {
	query := "SELECT * FROM participants WHERE id = ?"
	row := ps.db.QueryRow(query, participantID)

	var participant models.Participant
	err := row.Scan(
		&participant.ID,
		&participant.Name,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("participant not found with ID: %s", participantID)
	} else if err != nil {
		return nil, err
	}

	return &participant, nil
}
