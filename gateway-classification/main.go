package main

import (
	"github.com/stanleydv12/gateway-classification/src/database"
	"github.com/stanleydv12/gateway-classification/src/handler"
	"github.com/stanleydv12/gateway-classification/src/mqtt"
)

// main initializes the database and MQTT setup, then starts listening for MQTT messages.
//
// No parameters.
// No return types.
func main() {
	// Setup Database
	db := database.SetupDatabase()
	handler.SetDBInstance(db)

	// Setup Mqtt
	mqtt.SetupMqtt()

	// Application set to listen
	mqtt.Sub(mqtt.Client)

	handler.StartTimer()

	// Block the main function with a select statement
	select {}
}
