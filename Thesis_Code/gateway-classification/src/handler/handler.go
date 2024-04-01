package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/stanleydv12/gateway-classification/src/entity"
	"gorm.io/gorm"
)

var saveTimer *time.Timer

// HandleEvent handles the given event and data.
//
// Parameters:
// - event: a string representing the event to be handled.
// - data: an interface{} containing the data related to the event.
//
// Returns: None.
func HandleEvent(event string, data interface{}) {
	switch event {
	case "save-data":
		sensorData, ok := convertToSensorData(data)
		if !ok {
			fmt.Println("HandleEvent: Failed to convert data to SensorData")
			return
		}
		SaveData(sensorData)
	default:
		fmt.Println("HandleEvent: unknown event")
	}
}

// SaveData saves the provided data to the database if it meets certain conditions.
//
// The function takes in a data interface{} parameter, which should be of type entity.ECG.
// It first checks if the data is of the correct type using type assertion.
// If the data is not of type entity.ECG, the function prints an error message and returns.
//
// The function then parses the input timestamp of the data and checks if it is valid.
// If there is an error in parsing the timestamp, the function prints an error message and returns.
//
// The function calculates the time difference between the current time and the input timestamp.
// If the time difference is greater than 30 seconds, the function prints a message and returns.
//
// If the data is the first ECG for a patient, the function sets the ReferenceID to 0.
// Otherwise, it retrieves the ReferenceID for the patient using the GetReferenceIDForPatient function.
//
// The function creates a new instance of entity.ECG with the updated ReferenceID and other attributes.
//
// The function then attempts to save the newSensorData to the database using the DB.Create function.
// If there is an error in saving the data, the function prints an error message and returns.
//
// Finally, if the data is successfully saved to the database, the function prints a success message.
func SaveData(data interface{}) {
	fmt.Println(data)
	sensorData, ok := data.(entity.ECG)
	if !ok {
		fmt.Println("SaveData: Invalid data format")
		return
	}

	// Note: Inject input time with current time
	sensorData.InputTime = time.Now()
	currentTime := time.Now()

	diff := currentTime.Sub(sensorData.InputTime)
	if diff.Seconds() > 30 {
		fmt.Println("SaveData: Data is older than 30 seconds, skipping save to DB")
		return
	}

	if isFirstEcg := CheckIfFirstECG(sensorData); isFirstEcg {
		sensorData.ReferenceID = 0

		newSensorData := entity.ECG{
			ReferenceID: sensorData.ReferenceID,
			Value:       sensorData.Value,
			InputTime:   sensorData.InputTime,
		}

		err := DB.Create(&newSensorData).Error
		if err != nil {
			fmt.Println("SaveData: Failed to save data to DB:", err)
			return
		}

		firstSensorData := entity.ECG{}
		DB.Select("id").First(&firstSensorData)

		sleepData := entity.SleepData{
			ID:             1,
			PatientID:      1,
			FirstECGID:     firstSensorData.ID,
			FirstInputTime: firstSensorData.InputTime,
		}

		err = DB.Save(&sleepData).Error
		if err != nil {
			fmt.Println("SaveData: Failed to save sleep data to DB:", err)
			return
		}

	} else {
		sensorData.ReferenceID = GetReferenceIDForPatient(sensorData)

		newSensorData := entity.ECG{
			ReferenceID: sensorData.ReferenceID,
			Value:       sensorData.Value,
			InputTime:   sensorData.InputTime,
		}

		err := DB.Create(&newSensorData).Error
		if err != nil {
			fmt.Println("SaveData: Failed to save data to DB:", err)
			return
		}

		sleepData := entity.SleepData{
			ID:            1,
			PatientID:     1,
			LastInputTime: sensorData.InputTime,
		}

		err = DB.Save(&sleepData).Error
		if err != nil {
			fmt.Println("SaveData: Failed to save sleep data to DB:", err)
			return
		}
	}

	// Reset the timer to 30 seconds
	if saveTimer == nil {
		saveTimer = time.NewTimer(30 * time.Second)
	} else {
		saveTimer.Reset(30 * time.Second)
	}

	fmt.Println("SaveData: Data saved to DB")
}

// convertToSensorData converts the given data to a SensorData object.
//
// data: The data to be converted.
// Returns: The converted SensorData object and a boolean indicating if the conversion was successful.
func convertToSensorData(data interface{}) (entity.ECG, bool) {
	var sensorData entity.ECG
	err := mapstructure.Decode(data, &sensorData)
	if err != nil {
		return entity.ECG{}, false
	}
	return sensorData, true
}

var DB *gorm.DB

// SetDBInstance sets the global variable DB to the given *gorm.DB instance.
//
// db: the *gorm.DB instance to be set.
func SetDBInstance(db *gorm.DB) {
	DB = db
}

// CheckIfFirstECG checks if the given ECG sensor data is the first data recorded for a given date.
//
// It takes an ECG object as a parameter and returns a boolean value indicating whether the data is the first for that date.
func CheckIfFirstECG(sensorData entity.ECG) bool {
	var count int64
	var lastData entity.ECG

	// Check if data is empty
	DB.Model(&entity.ECG{}).Count(&count)
	if count == 0 {
		return true
	}

	// Check if the last data for reference_id is not empty
	DB.Order("input_time DESC").Where("reference_id = ?", sensorData.ReferenceID).First(&lastData)
	if lastData.ID != 0 {
		// Check if input_time is more than 30 seconds from the current time
		currentTime := time.Now()
		timeDiff := currentTime.Sub(lastData.InputTime)
		if timeDiff.Seconds() > 30 {
			return true
		}
	}

	return false
}

// GetReferenceIDForPatient retrieves the reference ID for a patient based on the given ECG sensor data.
//
// It takes the following parameter:
// - sensorData: an instance of the entity.ECG struct representing the ECG sensor data.
//
// It returns a uint representing the reference ID.
func GetReferenceIDForPatient(sensorData entity.ECG) uint {
	var referenceID uint
	DB.Model(&entity.ECG{}).Where("DATE(input_time) = ?", sensorData.InputTime.Format("2006-01-02")).Select("id").Order("id ASC").First(&referenceID)
	return referenceID
}

// StartTimer starts the timer for saving data to DB and calling classifyData
func StartTimer() {
	saveTimer = time.NewTimer(30 * time.Second)

	go func() {
		for {
			select {
			case <-saveTimer.C:
				classifyData()
				saveTimer.Reset(30 * time.Second)
			}
		}
	}()
}

type PredictionResponse struct {
	Prediction string `json:"prediction"`
}

func classifyData() {

	var lastECG entity.ECG
	if err := DB.Order("input_time desc").First(&lastECG).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("No ecg data found in the database")
		} else {
			log.Fatalf("Failed to get last ecg data: %v", err)
		}
		return
	}

	sleepData := entity.SleepData{
		ID:            1,
		LastInputTime: lastECG.InputTime,
	}
	if err := DB.Save(&sleepData).Error; err != nil {
		log.Fatalf("Failed to save sleep data: %v", err)
		return
	}
	

	// Mengambil data ECG berdasarkan ID
	var ecgByID []entity.ECG
	if err := DB.Where("id = ?", lastECG.ReferenceID).Find(&ecgByID).Error; err != nil {
		log.Fatalf("Failed to get ECG data by ID: %v", err)
		return
	}

	// Mengambil semua data ECG berdasarkan reference ID
	var ecgByReferenceID []entity.ECG
	if err := DB.Where("reference_id = ?", lastECG.ReferenceID).Find(&ecgByReferenceID).Error; err != nil {
		log.Fatalf("Failed to get ECG data by reference ID: %v", err)
		return
	}

	// Combine data retrieved by ID and reference ID
	allECG := append(ecgByID, ecgByReferenceID...)

	// Process data in batches of 10
	for i := 0; i < len(allECG); i += 10 {
		end := i + 10
		if end > len(allECG) {
			end = len(allECG)
		}

		batch := allECG[i:end]

		// Call the function to handle the API call for this batch
		prediction, err := predict(batch)
		if err != nil {
			log.Fatalf("Failed to make prediction: %v", err)
			return
		}

		if i == 0 {
			sleepStage := entity.SleepStage{
				ID:          1,
				ReferenceID: 0,
				Value:       prediction.Prediction,
			}

			err = DB.Save(&sleepStage).Error
			if err != nil {
				log.Fatalf("Failed to save sleep stage: %v", err)
				return
			}

			var sleepData entity.SleepData
			err = DB.Find(&sleepData, "id = ?", 1).First(&sleepData).Error
			if err != nil {
				log.Fatalf("Failed to get sleep data: %v", err)
				return
			}

			sleepData.FirstSleepStageID = sleepStage.ID
		} else {
			sleepStage := entity.SleepStage{
				ReferenceID: 1,
				Value:       prediction.Prediction,
			}

			err = DB.Save(&sleepStage).Error
			if err != nil {
				log.Fatalf("Failed to save sleep stage: %v", err)
				return
			}
		}
	}
}

func predict(data []entity.ECG) (PredictionResponse, error) {
	extractedValues := make([]float64, len(data))

	for i, ecg := range data {
		extractedValues[i] = ecg.Value
	}

	// Convert the extracted data slice to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return PredictionResponse{}, err
	}

	predictURL := os.Getenv("PREDICT_URL")
	if predictURL == "" {
		return PredictionResponse{}, fmt.Errorf("PREDICT_URL environment variable is not set")
	}

	resp, err := http.Post(predictURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return PredictionResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PredictionResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var predictionResponse PredictionResponse
	if err := json.NewDecoder(resp.Body).Decode(&predictionResponse); err != nil {
		return PredictionResponse{}, err
	}

	return predictionResponse, nil
}
