package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// DatabaseConnection represents the configuration for the database connection
type DatabaseConnection struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

// InitDB initializes a database connection
func InitDB(dbConfig DatabaseConnection) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Error pinging database: %v", err)
	}

	return db, nil
}
