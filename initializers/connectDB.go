package initializers

import (
	"fmt"
	"myapp/models"

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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)

	// Connect to the DB and initialize the DB variable
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("🚀 Connected Successfully to the Database")

	// DROP TABLES
	// DB.Migrator().DropTable(
	// 	&models.User{},
	// 	&models.UserProfile{},
	// 	// &models.Status{},
	// )

	// Migrate the database
	DB.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		// &models.Status{},
		// &models.Product{},
		// &models.Fact{},
	)

	fmt.Println("👍 Migration complete")

	// Initialize Status
	// var status = []models.Status{{Name: "Active"}, {Name: "Inactive"}, {Name: "Pending"}, {Name: "Suspended"}}
	// DB.Create(&status)

	SetUpDBConnection(DB)
}

func SetUpDBConnection(db *gorm.DB) {
	DB = db
}

func GetDBConnection() *gorm.DB {
	return DB
}
