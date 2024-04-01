package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Define your MQTT broker address and port
const (
	brokerAddress = "ws://localhost:9001" // Change this to your MQTT broker address
	host          = "localhost"
	port          = 5432
	user          = "postgres"
	password      = "postgres"
	dbname        = "sleep_monitoring"
)

var DB *sql.DB

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type SleepQuality struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	Value              string    `json:"value,omitempty"`
	DeepSleepDuration  float64   `json:"deep_sleep_duration,omitempty"`
	AwakeDuration      float64   `json:"awake_duration,omitempty"`
	TotalSleepDuration float64   `json:"total_sleep_duration,omitempty"`
	InputTime          time.Time `json:"input_time,omitempty"`
	UserId             uint      `json:"user_id,omitempty"`
}

// FuzzySet represents a fuzzy set with a membership function.
type FuzzySet struct {
	Name  string
	Terms map[string]func(float64) float64
}

// Fuzzify function calculates the membership values for each term in the fuzzy set.
func (fs *FuzzySet) Fuzzify(value float64) string {
	max := 0.0
	result := ""

	for term, mf := range fs.Terms {
		membership := mf(value)
		if membership > max {
			max = membership
			result = term
		}
	}

	return result
}

// ApplyRule function evaluates the conditions and returns the corresponding action.
func ApplyRule(conditions map[string]float64) (string, bool) {
	// Apply all rules
	switch {
	case conditions["awakeDuration"] == 1 && conditions["deepSleepTime"] == 1 && conditions["totalSleepTime"] == 1:
		return "LEVEL 9", true
	case conditions["awakeDuration"] == 1 && conditions["deepSleepTime"] == 1 && conditions["totalSleepTime"] == 0.5:
		return "LEVEL  8", true
	case conditions["awakeDuration"] == 1 && conditions["deepSleepTime"] == 1 && conditions["totalSleepTime"] == 0:
		return "LEVEL  7", true
	case conditions["awakeDuration"] == 1 && conditions["deepSleepTime"] == 0.5 && conditions["totalSleepTime"] == 1:
		return "LEVEL 7", true
	case conditions["awakeDuration"] == 1 && conditions["deepSleepTime"] == 0.5 && conditions["totalSleepTime"] == 0.5:
		return "LEVEL 6", true
	case conditions["awakeDuration"] == 1 && conditions["deepSleepTime"] == 0.5 && conditions["totalSleepTime"] == 0:
		return "LEVEL 5", true
	case conditions["awakeDuration"] == 1 && conditions["deepSleepTime"] == 0 && conditions["totalSleepTime"] == 1:
		return "LEVEL 5", true
	case conditions["awakeDuration"] == 1 && conditions["deepSleepTime"] == 0 && conditions["totalSleepTime"] == 0.5:
		return "LEVEL 4", true
	case conditions["awakeDuration"] == 1 && conditions["deepSleepTime"] == 0 && conditions["totalSleepTime"] == 0:
		return "LEVEL 3", true
	case conditions["awakeDuration"] == 0.5 && conditions["deepSleepTime"] == 1 && conditions["totalSleepTime"] == 1:
		return "LEVEL 8", true
	case conditions["awakeDuration"] == 0.5 && conditions["deepSleepTime"] == 1 && conditions["totalSleepTime"] == 0.5:
		return "LEVEL 7", true
	case conditions["awakeDuration"] == 0.5 && conditions["deepSleepTime"] == 1 && conditions["totalSleepTime"] == 0:
		return "LEVEL 6", true
	case conditions["awakeDuration"] == 0.5 && conditions["deepSleepTime"] == 0.5 && conditions["totalSleepTime"] == 1:
		return "LEVEL 6", true
	case conditions["awakeDuration"] == 0.5 && conditions["deepSleepTime"] == 0.5 && conditions["totalSleepTime"] == 0.5:
		return "LEVEL 5", true
	case conditions["awakeDuration"] == 0.5 && conditions["deepSleepTime"] == 0.5 && conditions["totalSleepTime"] == 0:
		return "LEVEL 4", true
	case conditions["awakeDuration"] == 0.5 && conditions["deepSleepTime"] == 0 && conditions["totalSleepTime"] == 1:
		return "LEVEL 4", true
	case conditions["awakeDuration"] == 0.5 && conditions["deepSleepTime"] == 0 && conditions["totalSleepTime"] == 0.5:
		return "LEVEL 3", true
	case conditions["awakeDuration"] == 0.5 && conditions["deepSleepTime"] == 0 && conditions["totalSleepTime"] == 0:
		return "LEVEL 2", true
	case conditions["awakeDuration"] == 0 && conditions["deepSleepTime"] == 1 && conditions["totalSleepTime"] == 1:
		return "LEVEL 7", true
	case conditions["awakeDuration"] == 0 && conditions["deepSleepTime"] == 1 && conditions["totalSleepTime"] == 0.5:
		return "LEVEL 6", true
	case conditions["awakeDuration"] == 0 && conditions["deepSleepTime"] == 1 && conditions["totalSleepTime"] == 0:
		return "LEVEL 5", true
	case conditions["awakeDuration"] == 0 && conditions["deepSleepTime"] == 0.5 && conditions["totalSleepTime"] == 1:
		return "LEVEL 5", true
	case conditions["awakeDuration"] == 0 && conditions["deepSleepTime"] == 0.5 && conditions["totalSleepTime"] == 0.5:
		return "LEVEL 4", true
	case conditions["awakeDuration"] == 0 && conditions["deepSleepTime"] == 0.5 && conditions["totalSleepTime"] == 0:
		return "LEVEL 3", true
	case conditions["awakeDuration"] == 0 && conditions["deepSleepTime"] == 0 && conditions["totalSleepTime"] == 1:
		return "LEVEL 3", true
	case conditions["awakeDuration"] == 0 && conditions["deepSleepTime"] == 0 && conditions["totalSleepTime"] == 0.5:
		return "LEVEL 2", true
	case conditions["awakeDuration"] == 0 && conditions["deepSleepTime"] == 0 && conditions["totalSleepTime"] == 0:
		return "LEVEL 1", true
	}

	// No rule matched
	return "", false
}

var messagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("Message Received from [%s] : %s\n", msg.Topic(), msg.Payload())

	if string(msg.Payload()) == "sleepStageReady" {
		quantifyData()
	}
}

func setUpDB() *sql.DB {
	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	// Attempt to ping the database to verify connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}

	fmt.Println("Successfully connected to the database!")

	return db
}

var red *redis.Client

func quantifyData() {

	var awake float64 = 0
	var deepSleep float64 = 0
	var totalSleep float64 = 0

	db := setUpDB()
	DB = db
	value, err := red.Get(red.Context(), "sleep-stages").Result()
	if value == "" || err != nil {
		// Example query
		start := time.Now().UnixMilli()
		rows, err := DB.Query("SELECT * FROM sleep_stages")
		end := time.Now().UnixMilli()
		fmt.Println("Start : ", start)
		fmt.Println("End : ", end)
		fmt.Println("Duration : ", end-start, "ms")
		if err != nil {
			log.Fatal("Error querying the database: ", err)
		}
		defer rows.Close()

		// Iterate over the rows
		for rows.Next() {
			var id int
			var reference_id string
			var value string
			var method any
			var timestamp time.Time
			err := rows.Scan(&id, &reference_id, &value, &method, &timestamp)
			if err != nil {
				log.Println("Error scanning row: ", err)
			}
			if value == "N1" || value == "N2" || value == "REM" {
				totalSleep++
			}
			if value == "N3" {
				totalSleep++
				deepSleep++
			}
			if value == "AWAKE" {
				awake++
			}
		}
		if err := rows.Err(); err != nil {
			log.Println("Error iterating over rows: ", err)
		}
	} else {
		type SleepStage struct {
			ID          uint      `json:"id"`
			ReferenceID uint      `json:"reference_id"`
			Value       string    `json:"value"`
			Method      string    `json:"method"`
			Timestamp   time.Time `json:"timestamp"`
		}

		type SleepStages struct {
			Stages []SleepStage `json:"sleep-stages"`
		}

		// Unmarshal the JSON into a SleepStages struct
		var stages SleepStages
		err := json.Unmarshal([]byte(value), &stages)
		if err != nil {
			log.Println(err)
		}

		// Loop through the stages and print values
		for _, stage := range stages.Stages {
			if stage.Value == "N1" || stage.Value == "N2" || stage.Value == "REM" {
				totalSleep++
			}
			if stage.Value == "N3" {
				totalSleep++
				deepSleep++
			}
			if stage.Value == "AWAKE" {
				awake++
			}
		}
	}

	// Define fuzzy sets and input values
	awakeDurationSet := FuzzySet{
		Name: "awakeDuration",
		Terms: map[string]func(float64) float64{
			"low":    func(x float64) float64 { return math.Max(0, 1.0-0.5*x) },
			"medium": func(x float64) float64 { return math.Max(0, 1-2*math.Abs(x-0.5)) },
			"high":   func(x float64) float64 { return math.Max(0, 2*x-1) },
		},
	}
	deepSleepDurationSet := FuzzySet{
		Name: "awakeDuration",
		Terms: map[string]func(float64) float64{
			"low":    func(x float64) float64 { return math.Max(0, 1.0-0.5*x) },
			"medium": func(x float64) float64 { return math.Max(0, 1-2*math.Abs(x-0.5)) },
			"high":   func(x float64) float64 { return math.Max(0, 2*x-1) },
		},
	}
	totalSleepDurationSet := FuzzySet{
		Name: "awakeDuration",
		Terms: map[string]func(float64) float64{
			"low":    func(x float64) float64 { return math.Max(0, 1.0-0.5*x) },
			"medium": func(x float64) float64 { return math.Max(0, 1-2*math.Abs(x-0.5)) },
			"high":   func(x float64) float64 { return math.Max(0, 2*x-1) },
		},
	}

	inputValues := map[string]float64{
		"awakeDuration":  awake / 60,
		"deepSleepTime":  deepSleep / 60,
		"totalSleepTime": totalSleep / 60,
	}

	// Fuzzify the sleep data
	awakeTerm := awakeDurationSet.Fuzzify(inputValues["awakeDuration"])
	deepSleepTerm := deepSleepDurationSet.Fuzzify(inputValues["deepSleepTime"])
	totalSleepTerm := totalSleepDurationSet.Fuzzify(inputValues["totalSleepTime"])
	// Fuzzify the other sleep data...

	// Print the intermediate results
	fmt.Println("Fuzzy Conditions:")
	fmt.Printf("Awake Duration: %s (%f)\n", awakeTerm, inputValues["awakeDuration"])
	fmt.Printf("Deep Sleep Duration: %s (%f)\n", deepSleepTerm, inputValues["deepSleepTime"])
	fmt.Printf("Total Sleep Duration: %s (%f)\n", totalSleepTerm, inputValues["totalSleepTime"])
	// Print the other fuzzy conditions...

	// Map the linguistic terms to their respective numerical values
	linguisticValues := map[string]float64{
		"high":   1.0,
		"medium": 0.5,
		"low":    0.0,
	}

	// Create a map for all conditions
	allConditions := map[string]float64{
		"awakeDuration":  linguisticValues[awakeTerm],
		"deepSleepTime":  linguisticValues[deepSleepTerm],
		"totalSleepTime": linguisticValues[totalSleepTerm],
	}

	action, ruleMatched := ApplyRule(allConditions)

	// Print the final results
	if ruleMatched {
		fmt.Println("\nRule Matched! Action:", action)
	} else {
		fmt.Println("\nNo Rule Matched.")
	}

	sleepQuality := SleepQuality{
		Value:              action,
		DeepSleepDuration:  inputValues["deepSleepTime"],
		TotalSleepDuration: inputValues["totalSleepTime"],
		AwakeDuration:      inputValues["awakeDuration"],
		InputTime:          time.Now(),
		UserId:             1,
	}

	sleepQualityJSON, err := json.Marshal(sleepQuality)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = red.Set(red.Context(), "sleep_quality_1", sleepQualityJSON, 0).Err()
	if err != nil {
		log.Fatalf("Failed to save to redis: %v", err)
	} else {
		fmt.Println("Saved to redis")
	}

	// Prepare the SQL statement for inserting data.
	stmt, err := db.Prepare("INSERT INTO sleep_qualities (value, deep_sleep_duration, awake_duration, total_sleep_duration, input_time, user_id) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Example data to be inserted.
	value1 := action
	value2 := inputValues["deepSleepTime"]
	value3 := inputValues["awakeDuration"]
	value4 := inputValues["totalSleepTime"]
	value5 := time.Now()
	value6 := 1

	// Execute the SQL statement with the provided values.
	_, err = stmt.Exec(value1, value2, value3, value4, value5, value6)
	if err != nil {
		log.Println(err)
	}

}

func main() {
	red = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"), // Redis server address
		Password: "",                      // no password set
		DB:       0,                       // use default DB
	})

	// Set MQTT client options

	opts := MQTT.NewClientOptions().AddBroker("ws://" + os.Getenv("MQTT_BROKER"))
	opts.SetClientID("MQTT_simulator")
	opts.SetDefaultPublishHandler(messagePubHandler)

	topic := "sleep_monitoring"

	// Create and start a client using the above options
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Error connecting to MQTT broker:", token.Error())
		os.Exit(1)
	}

	// Subscribe to MQTT topics
	if token := client.Subscribe(topic, 0, messagePubHandler); token.Wait() && token.Error() != nil {
		fmt.Println("Error subscribing to MQTT topic:", token.Error())
		os.Exit(1)
	}
	fmt.Printf("Subscribed to topic: %s\n", topic)
	// Keep the program running indefinitely
	select {}
}
