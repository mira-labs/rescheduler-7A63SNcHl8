// Package main is the entry point for the rescheduler application.
// It contains the main function to demonstrate the usage of the rescheduler logic
// and interacts with the LambdaHandler for handling questionnaire completion events.
//
// The rescheduler application is designed to process questionnaire completion events,
// update the database, and handle scheduling logic for new questionnaires.
// It uses various internal packages for database operations, SQS messaging, and utility functions.
// Because the setup for lambda integration is not known, certain assumptions have beenn made:
package main

import (
	"context"
	"fmt"
	"time"

	"rescheduler/internals/database"
	"rescheduler/internals/models"
	"rescheduler/internals/sqs"
	"rescheduler/internals/store"
	"rescheduler/internals/timestamp"
	"rescheduler/internals/util"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"umotif.com/go/credentials"
)

type DatabaseConnection struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

func main() {

	//credentials.MongoDbDsn() could provide the same result as the below one

	//Example event json. The mapping template and the integrtion settings between API Gateway and the Lambda function is not known at this stage.
	jsonString := `{
		"id": "random",
		"user_id": "8a4378cd-27b9-4a36-afee-829b42eeb1b5",
		"study_id": "Study5",
		"questionnaire_id": "24b6f062-df29-4e6a-abb4-403e01671e4a",
		"completed_at": "2023-12-04 02:11:00",
		"remaining_completions": 2
	}`
	fmt.Printf("Event data: %s", jsonString+"\n")

	// Creating a variable of type models.QuestionnaireCompletedEvent to handle event json
	var eventData *models.QuestionnaireCompletedEvent
	var err error

	eventData, err = util.ConvertJSONToEvent(jsonString)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Printf("Questionnaire: %s for parrticipant: %s\n", eventData.UserID, eventData.QuestionnaireID)
	fmt.Println("Completed at: ", eventData.CompletedAt)
	// Invoke the LambdaHandler
	lambda.Start(LambdaHandler)

}

// LambdaHandler is the AWS Lambda handler function for processing questionnaire completion events.
// It connects to the database, retrieves information about the completed questionnaire and schedule,
// updates the schedule status, creates a new schedule if needed, sends SQS messages, and records the questionnaire result.
//
// The function uses channels for synchronization to find the questionnaire and schedule simultaneously.
// Once both tasks are completed, it updates the schedule status asynchronously, creates a new schedule if needed,
// sends SQS messages, and records the questionnaire result concurrently.
//
// Parameters:
//   - ctx: A context.Context object.
//   - event: A pointer to the models.QuestionnaireCompletedEvent containing the event data.
//
// Returns:
//   - An events.APIGatewayProxyResponse indicating the success or failure of the operation.
//   - An error if there was an internal server error.
func LambdaHandler(ctx context.Context, event *models.QuestionnaireCompletedEvent) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Connecting to database")

	//Create an SQS handler instance
	sqsHandler := sqs.NewSQSHandler("sqs_url", "aws_region")

	//Connect to database
	dbConfig := credentials.MySQLDbDsn()
	db, err := database.InitDB(dbConfig)
	if err != nil {
		fmt.Println("Error: ", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Internal server error"),
			StatusCode: 500,
		}, err
	}

	fmt.Println("Connected to db: ", dbConfig.Database)
	defer db.Close()

	// Create instances of the database stores
	questionnaireStore := store.NewQuestionnaireStore(db)
	scheduledQuestionnaireStore := store.NewScheduledQuestionnaireStore(db)
	questionnaireResultStore := store.NewQuestionnaireResultStore(db)

	// Use channels for synchronization
	resultCh := make(chan struct {
		*models.Questionnaire
		*models.ScheduledQuestionnaire
	})
	errorCh := make(chan error)

	// Find the schedule and questionnaire simultaneously
	go func() {
		questionnaire, err := questionnaireStore.FindQuestionnaireByID(event.QuestionnaireID)
		if err != nil {
			errorCh <- err
			return
		}

		schedule, err := scheduledQuestionnaireStore.FindScheduledQuestionnaireByQuestionnaireIDAndUserID(
			event.QuestionnaireID, event.UserID,
		)
		if err != nil {
			errorCh <- err
			return
		}

		resultCh <- struct {
			*models.Questionnaire
			*models.ScheduledQuestionnaire
		}{questionnaire, schedule}
	}()

	// Wait for the goroutine to complete
	select {
	case result := <-resultCh:
		// Both tasks completed successfully
		questionnaire := result.Questionnaire
		schedule := result.ScheduledQuestionnaire

		// Update the schedule status to 'completed' asynchronously
		go func() {
			err := scheduledQuestionnaireStore.SetScheduledQuestionnaireCompleted(schedule)
			if err != nil {
				fmt.Println("Error updating schedule: ", err)
			}
		}()

		//Checking if there are remaining completions or if the max_attempt in the database is NULL
		if event.RemainingCompletions > 0 || !questionnaire.MaxAttempts.Valid {
			go func() {

				// Create a new schedule for the same questionnaire
				nextScheduledTime := timestamp.TimeStamp{
					Time: event.CompletedAt.Add(
						time.Duration(questionnaire.HoursBetweenAttempts) * time.Hour,
					),
				}

				// Debugging output
				fmt.Println("Original CompletedAt:", event.CompletedAt)
				fmt.Println("HoursBetweenAttempts:", questionnaire.HoursBetweenAttempts)

				// Debugging output
				fmt.Println("Next Scheduled Time:", nextScheduledTime)

				newScheduledQuestionnaire := models.ScheduledQuestionnaire{
					ID:              uuid.New().String(),
					QuestionnaireID: questionnaire.ID,
					ParticipantID:   event.UserID,
					ScheduledAt:     nextScheduledTime,
					Status:          models.ScheduledQuestionnairePending,
				}
				fmt.Println("Created a new Scheduled Questionnaire:", newScheduledQuestionnaire.ID)

				// Save the new schedule
				err := scheduledQuestionnaireStore.Create(&newScheduledQuestionnaire)
				if err != nil {
					fmt.Println("Error: ", err)
				}

				fmt.Println("Saved the Scheduled Questionnaire: ", newScheduledQuestionnaire.ID)

				sqsHandler.SendNewScheduleMessage(newScheduledQuestionnaire.ID, event.UserID)

			}()

		} else {
			go func() {
				sqsHandler.SendCompletionMessage(event.UserID)
			}()
		}

		// Create a questionnaire result asynchronously
		go func() {
			questionnaireResult := models.QuestionnaireResult{
				ID:                      uuid.New().String(),
				Answers:                 `{"question":"answer"}`,
				QuestionnaireID:         questionnaire.ID,
				ParticipantID:           event.UserID,
				QuestionnaireScheduleID: schedule.ID,
				CompletedAt:             event.CompletedAt,
			}

			err := questionnaireResultStore.Create(&questionnaireResult)
			if err != nil {
				fmt.Println("Error creating questionnaire result: ", err)
			}
		}()

		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Success"),
			StatusCode: 200,
		}, nil

	case err := <-errorCh:
		// Handle any error that occurred during finding schedule or questionnaire
		fmt.Println("Error: ", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Internal server error"),
			StatusCode: 500,
		}, nil
	}
}
