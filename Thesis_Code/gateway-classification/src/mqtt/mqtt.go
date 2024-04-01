package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
	"github.com/stanleydv12/gateway-classification/src/handler"
)

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

var Client mqtt.Client

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Message Received from [%s] : %s\n", msg.Topic(), msg.Payload())

	var receivedMessage Message
	err := json.Unmarshal(msg.Payload(), &receivedMessage)
	if err != nil {
		fmt.Println("Error decoding message:", err)
		return
	}

	handler.HandleEvent(receivedMessage.Event, receivedMessage.Data)
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
	Sub(client)
}

var onLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost : %v", err)
}

// SetupMqtt sets up the MQTT connection.
//
// It loads the environment variables from the .env file.
// It creates a new MQTT client options object and configures it with the
// MQTT broker host, port, and client ID.
// It sets the default publish handler, on connect handler, and connection lost
// handler for the MQTT client options.
// It creates a new MQTT client with the configured options.
// It connects the MQTT client to the broker.
// If there is an error during the connection, it panics.
func SetupMqtt() {
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Fatal("Error loading .env")
	}

	broker := os.Getenv("MQTT_BROKER_HOST")
	port := os.Getenv("MQTT_BROKER_PORT")
	clientID := os.Getenv("MQTT_CLIENT_ID")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ws://%s:%s", broker, port))
	opts.SetClientID(clientID)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetOnConnectHandler(connectHandler)
	opts.SetConnectionLostHandler(onLostHandler)

	Client = mqtt.NewClient(opts)

	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

// Sub is a function that subscribes to an MQTT topic and prints a message.
//
// It takes a client of type mqtt.Client as a parameter.
// It does not return any value.
func Sub(client mqtt.Client) {
	topic := os.Getenv("MQTT_TOPIC")
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s\n", topic)
}

// Pub publishes the given message to the MQTT topic.
//
// It takes a single parameter, msg, which is of type interface{}.
// The function does not return any value.
func Pub(msg interface{}) {
	topic := os.Getenv("MQTT_TOPIC")
	token := Client.Publish(topic, 1, false, msg)
	token.Wait()
}
