package configs

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Declare the variable for the database
var DB *gorm.DB

// ConnectDB connect to db
func ConnectDB(config *Config) {
	var err error
	// Connection URL to connect to Postgres Database
	dsn := config.DB_DSN

	loggerLevel := logger.Info
	if !config.IsDebug {
		loggerLevel = logger.Warn
	}

	// Connect to the DB and initialize the DB variable
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(loggerLevel),
	})

	// enable UUID support
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("🚀 Connected Successfully to the Database")

	setUpDBConnection(DB)
}

func setUpDBConnection(db *gorm.DB) {
	DB = db
}

func GetDBConnection() *gorm.DB {
	return DB
}

func MigrateDB() {
	// DROP TABLES
	// DB.Migrator().DropTable(
	// 	// &models.User{},
	// 	// &models.UserProfile{},
	// 	// &models.OTPRequest{},
	// 	// &models.Product{},
	// 	&models.MyDrive{},
	// )

	// Migrate the database
	// DB.AutoMigrate(
	// 	// &models.User{},
	// 	// &models.UserProfile{},
	// 	// &models.OTPRequest{},
	// 	// &models.Product{},
	// 	&models.MyDrive{},
	// )

	// fmt.Println("👍 Migration complete")

	// Initialize Status
	// DB.AutoMigrate(&models.Status{})
	// var status = []models.Status{{Name: "Active"}, {Name: "Inactive"}, {Name: "Pending"}, {Name: "Suspended"}}
	// DB.Create(&status)
}
