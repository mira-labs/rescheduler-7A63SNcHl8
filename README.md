# Rescheduler

The rescheduler application is designed to process questionnaire completion events,
update the database, and handle scheduling logic for new questionnaires.
It uses various internal packages for database operations, SQS messaging, and utility functions.

## Assumptions

Because the setup for lambda integration is not known, certain assumptions have been made:

* Lambda function returns a `APIGatewayProxyResponse` HTTP response through `github.com/aws/aws-lambda-go/events`
* The `credentials.MySQLDbDsn()` method from `umotif.com/go/credentials` package is used for obtaining database connection credentials.
* The SQS URL and AWS region are hardcoded for demonstration purposes.
* The `StudyID` field, even though present in the event data, does not serve as an identifier for the questionnaire. As a result, it is not required to be utilized in the lookup functions.
* The application does not manage the count of remaining attempts since this information is included in the event data.
* This implementation does not use GORM or any other ORM-like package or framework as the number of models and database operations is tiny.

## Installation

```bash

# Building the binary
make build

# Cleaning up
make clean

# Running
./bin/main
```

## Architecture

I have taken an __Onion__ approach following some of the Domain-Driven Design (DDD) principles separating the domain, presentation, application, and infrastructure layers. 

The project consists of separate files for handling database operations [`internals/database/database.go`](internals/database/database.go), interacting with data stores, and handling Simple Queue Service (SQS) messaging ([`internals/sqs/sqs.go`](./internals/sqs/sqs.go).

### Package `models`

The `models` package in the application serves as a pivotal component, providing structured data representations for key entities involved in the study scheduling and completion process. These entities include participants, questionnaires, scheduled questionnaires, and questionnaire results.

They represent the following database structure:

![Database structure](assets/db-layout-diagram.png)

The `Participant` structure represents an individual participating in a study, characterized by a unique identifier (`ID`) and a name. 

Questionnaires, the core instruments for collecting study data, are modeled by the `Questionnaire` structure. This structure includes details such as the questionnaire's ID, associated study ID, name, question configurations, maximum attempts, and scheduling parameters.

Scheduled instances of questionnaires are captured by the `ScheduledQuestionnaire` structure. This structure includes information such as the scheduled questionnaire's ID, associated questionnaire and participant IDs, scheduled timestamp, and status. Finally, the `QuestionnaireResult` structure encapsulates the results of a participant completing a questionnaire, storing data such as result ID, answers, associated questionnaire and participant IDs, the schedule ID, and completion timestamp.

In addition to these core structures, the `models` package features a `QuestionnaireCompletedEvent` structure. This model is designed to represent a specific event related to questionnaire completion. It includes properties such as event ID, user ID, study ID, questionnaire ID, completion timestamp, and the count of remaining completions. The overall structure of the `models` package establishes a robust foundation for database interactions, ensuring organized and standardized representations of study-related entities in the Rescheduler application. These structures are instrumental in maintaining the integrity and coherence of the application's data model throughout various operations.

### Package `database`

#### [`internals/database/database.go`](./internals/database/database.go)

The `database.go` file defines a `DatabaseConnection` struct for configuring the database connection and an `InitDB` function for initializing a connection to a MySQL database. This file encapsulates the logic for establishing and validating the database connection, providing a clean and centralized way to manage database interactions throughout the application.


### Package `store`


* [`internals/store/questionnaire_store.go`](internals/store/questionnaire_store.go)
* [`internals/store/cheduled_questionnaire_store.go`](internals/store/cheduled_questionnaire_store.go)
* [`internals/store/participant_store.go`](internals/store/participant_store.go)
* [`internals/store/uestionnaire_result_store.go`](internals/store/uestionnaire_result_store.go)
 
These files are responsible for interacting with specific entities in the database. Each file defines a corresponding store type (`QuestionnaireStore`, `ScheduledQuestionnaireStore`, `ParticipantStore`, `QuestionnaireResultStore`) that encapsulate database operations for its respective entity. This modular approach adheres to the Single Responsibility Principle, making it easier to maintain and extend the codebase.

### Package `sqs`

#### [`internals/sqs/sqs.go`](./internals/sqs/sqs.go)

This file handles SQS messaging, providing an abstraction for sending messages to an SQS queue. It defines a `SQSHandler` type that encapsulates the logic for sending completion and new schedule messages to the SQS queue. 


### Package `timestamp`

#### [`internals/timestamp/timestamp.go`](./internals/timestamp/timestamp.go)

The `timestamp.go` file defines a custom time type named `TimeStamp`, which extends the standard `time.Time` functionality. The primary purpose of this file is to facilitate the integration of MySQL `DATETIME` values with the `time.Time` type used in the models. 

The `Scan` method implements the `sql.Scanner` interface. This allows instances of `TimeStamp` to be seamlessly scanned from database query results. It converts a raw database value (in this case, a MySQL `DATETIME` value) into a `time.Time` value and sets it as the underlying `time.Time` field of the `TimeStamp` type. 

The `UnmarshalJSON` method implements JSON unmarshalling specifically for the `TimeStamp` type. It unmarshals a JSON byte slice into a string, and then it parses that string into a `time.Time` value using the defined layout. This method is useful when dealing with JSON data, ensuring that the custom time type can be unmarshalled correctly.

### Package `util`

#### [`internals/util/util.go`](./internals/util/util.go)

The util.go file provides utility functions for the application. It specifically contains the `ConvertJSONToEvent` function, which facilitates the conversion of a JSON string into a `models.QuestionnaireCompletedEvent` struct. This utility is essential for handling incoming JSON data, such as questionnaire completion events, and converting it into a structured format that can be used within the application.


## Rescheduler Application Logic Overview

1. Example JSON is parsed and converted to `QuestionnaireCompletedEvent`
2. The object is then passed to the `LambdaHandler` function.
3. Database Connection and SQS service handler instance:
   
   The function starts by establishing a connection to the database using the provided database configuration.
    ```go
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
    ```
4. Database Stores
   
   Instances of the database stores (`questionnaireStore`, `scheduledQuestionnaireStore`, `questionnaireResultStore`) are created using the initialized database connection.
   `ParticipantStore` despite being implemented is not used because in this example there are no operations on participants in the business logic layer.
   
    ```go
    questionnaireStore := store.NewQuestionnaireStore(db)
    scheduledQuestionnaireStore := store.NewScheduledQuestionnaireStore(db)
    questionnaireResultStore := store.NewQuestionnaireResultStore(db)
    ```
6. Channel Initialization

    Two channels are created for synchronization:
       * `resultCh`: To receive the results of finding the questionnaire and schedule.
       * `errorCh`: To receive any errors that might occur during the database operations.
    
    ```go
    resultCh := make(chan struct {
    *models.Questionnaire
    *models.ScheduledQuestionnaire
    })
    errorCh := make(chan error)
     ```
7. Concurrent Tasks

    A goroutine is spawned to concurrently find both the questionnaire and the schedule.
    If an error occurs during any of these tasks, the error is sent to the `errorCh`. If both tasks are successful, the results are sent to the `resultCh`.

    ```go
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
    ```
8. Channel Selection
   
   The select statement is used to wait for either the successful completion of finding both the questionnaire and schedule `(resultCh)` or an error `(errorCh)`.
   Also, the first operation is to set the questionnaire status to `complete`.

   ```go
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

9. Checking if a new schedule can be created
   * Are there any more attempts left or is there no limit (`max_attempts` field in the database is NULL)?
   * If that's the case, create a new schedule `questionnaire.HoursBetweenAttempts` after `event.CompletedAt`
   * Send SQS message confirming creation.
   * If not, a "completion" SQS message is sent.
     
   I'm assuming the SQS messages have to be sent only after the database events are completed, therefore no goroutines are used between them.

   ```go
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
					Status:          "pending",
				}
				fmt.Println("Created a new Scheduled Questionnaire:", newScheduledQuestionnaire.ID)

				// Save the new schedule
				err := scheduledQuestionnaireStore.CreateScheduledQuestionnaire(&newScheduledQuestionnaire)
				if err != nil {
					fmt.Println("Error: ", err)
				}

				fmt.Println("Saved the Scheduled Questionnaire: ", newScheduledQuestionnaire.ID)

				sqsHandler.SendNewScheduleMessage(newScheduledQuestionnaire.ID, event.UserID)

			}()

		} else {
			sqsHandler.SendCompletionMessage(event.UserID)
		}
   ```

10. Create a questionnaire result record.
   This is done asynchronously as this happens regardless of the results of the previous step.

     ```go
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

                err := questionnaireResultStore.CreateQuestionnaireResult(&questionnaireResult)
                if err != nil {
                    fmt.Println("Error creating questionnaire result: ", err)
                }
            }()
            ```
11. Return the successful response
    ```go
            return events.APIGatewayProxyResponse{
                Body:       fmt.Sprintf("Success"),
                StatusCode: 200,
            }, nil
    ```

12. Handle the top-level errors that occurred during finding the schedule or questionnaire
     
    ```go
        case err := <-errorCh:
            fmt.Println("Error: ", err)
            return events.APIGatewayProxyResponse{
                Body:       fmt.Sprintf("Internal server error"),
                StatusCode: 500,
            }, nil
        }
    ```




