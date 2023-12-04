USE `scheduled_questionnaires`;

CREATE TABLE IF NOT EXISTS participants (
    id VARCHAR(128) PRIMARY KEY NOT NULL,
    name VARCHAR(128) NOT NULL
);

CREATE TABLE IF NOT EXISTS questionnaires (
    id VARCHAR(128) PRIMARY KEY NOT NULL,
    study_id VARCHAR(128) NOT NULL,
    name VARCHAR(128) NOT NULL,
    questions JSON NOT NULL,
    max_attempts INT,
    hours_between_attempts INT DEFAULT 24
);

CREATE TABLE IF NOT EXISTS scheduled_questionnaires (
    id VARCHAR(128) PRIMARY KEY NOT NULL,
    questionnaire_id VARCHAR(128) NOT NULL,
    participant_id VARCHAR(128) NOT NULL,
    scheduled_at DATETIME NOT NULL,
    status ENUM('pending', 'completed') NOT NULL
);

CREATE TABLE IF NOT EXISTS questionnaire_results (
    id VARCHAR(128) NOT NULL,
    answers JSON NOT NULL,
    questionnaire_id VARCHAR(128) NOT NULL,
    participant_id VARCHAR(128) NOT NULL,
    questionnaire_schedule_id VARCHAR(128),
    completed_at DATETIME
);