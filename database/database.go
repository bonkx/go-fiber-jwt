package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

// GlobalDB is a global db object that will be used across different packages
var Database *gorm.DB

// InitDatabase creates a mysql db connection and stores it in the GlobalDB variable
// It reads the environment variables from the .env file and uses them to create the connection
// It returns an error if the connection fails
func InitDatabase() {
	var err error

	// Read the environment variables from the .env file
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	databaseName := os.Getenv("DB_NAME")

	// Create the data source name (DSN) using the environment variables
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, databaseName,
	)
	// Create the connection and store it in the GlobalDB variable
	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully connected to the database")
	}
}
