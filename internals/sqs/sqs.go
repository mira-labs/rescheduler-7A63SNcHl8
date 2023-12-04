// Package sqs provides an interface and implementation for interacting with Amazon Simple Queue Service (SQS).
// It includes methods for sending messages related to schedule creation and completion.
package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQS defines methods for interacting with Amazon SQS.
type SQS interface {
	SendNewScheduleMessage(scheduleID string, participantID string) error
	SendCompletionMessage(participantID string) error
}

// SQSHandler is an implementation of the SQS interface.
type SQSHandler struct {
	queueURL string
	region   string
}

// NewSQSHandler creates a new instance of SQSImpl.
func NewSQSHandler(queueURL, region string) *SQSHandler {
	return &SQSHandler{
		queueURL: queueURL,
		region:   region,
	}
}

// SendNewScheduleMessage sends a message to SQS indicating that a new schedule has been created.
// The actual message structure is not known at this stage so a simple text message is sent.
func (s *SQSHandler) SendNewScheduleMessage(scheduleID string, participantID string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s.region),
	})
	if err != nil {
		return err
	}

	svc := sqs.New(sess)

	message := "New schedule created with ID: " + scheduleID + " for participant with ID: " + participantID

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		MessageBody:  aws.String(message),
		QueueUrl:     &s.queueURL,
		DelaySeconds: aws.Int64(0),
	})

	return err
}

// SendCompletionMessage sends a message to SQS indicating that the user has completed all scheduled questionnaires.
// The actual message structure is not known at this stage, so a simple text message is sent.
func (s *SQSHandler) SendCompletionMessage(participantID string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s.region),
	})
	if err != nil {
		return err
	}

	svc := sqs.New(sess)

	message := "User " + participantID + " has completed all scheduled questionnaires."

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		MessageBody:  aws.String(message),
		QueueUrl:     &s.queueURL,
		DelaySeconds: aws.Int64(0),
	})

	return err
}
