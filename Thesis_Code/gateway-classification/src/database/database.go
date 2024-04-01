package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/stanleydv12/gateway-classification/src/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupDatabase initializes and returns a *gorm.DB object representing a connection to the database.
//
// It reads the database configuration from the environment variables and establishes a connection to the database using the specified parameters.
// The function automatically migrates the entity.SensorData struct to the database and returns the *gorm.DB object.
func SetupDatabase() *gorm.DB {

	var err error

	errEnv := godotenv.Load(".env")

	if errEnv != nil {
		log.Fatal("Error loading .env")
	}

	host := os.Getenv("DSN_HOST")
	user := os.Getenv("DSN_USER")
	password := os.Getenv("DSN_PASSWORD")
	dbname := os.Getenv("DSN_DB_NAME")
	port := os.Getenv("DSN_PORT")
	sslMode := os.Getenv("SSL_MODE")
	timeZone := os.Getenv("TIMEZONE")

	dsn :=
		"host=" + host +
			" user=" + user +
			" password=" + password +
			" dbname=" + dbname +
			" port=" + port +
			" sslmode=" + sslMode +
			" TimeZone=" + timeZone

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	err = AutoMigrateAllEntities(db)
	if err != nil {
		log.Fatalf("Error migrating entities: %v", err)
	}

	// Seeding
	InsertTestData(db)

	return db
}

// AutoMigrateAllEntities migrates all entities to the database.
func AutoMigrateAllEntities(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entity.Doctor{},
		&entity.ECG{},
		&entity.Patient{},
		&entity.SleepData{},
		&entity.SleepStage{},
		&entity.SleepQuality{},
	)
	if err != nil {
		return err
	}
	return nil
}

func InsertTestData(db *gorm.DB) {
	// Contoh data Doctor
	doctorData := entity.Doctor{
		FirstName:        "John",
		LastName:         "Doe",
		Gender:           "Male",
		Email:            "john.doe@example.com",
		PhoneNumber:      "123456789",
		BirthRate:        time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
		RegistrationDate: time.Now(),
		Address:          "123 Main St, City",
	}

	// Contoh data Patient
	patientData := entity.Patient{
		FirstName:        "Alice",
		LastName:         "Johnson",
		Gender:           "Female",
		Email:            "alice.johnson@example.com",
		PhoneNumber:      "987654321",
		BirthDate:        time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
		RegistrationDate: time.Now(),
		Address:          "456 Oak St, Town",
		SensorToken:      "sensortoken123",
		StartSleepTime:   time.Now(),
		EndSleepTime:     time.Now().Add(8 * time.Hour), // Contoh: tidur selama 8 jam
	}

	// Sisipkan data ke database
	db.Create(&doctorData)
	db.Create(&patientData)

	fmt.Println("Data inserted successfully.")
}
