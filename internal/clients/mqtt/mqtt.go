package mqtt

import (
	"encoding/json"
	"fmt"
	mqtt2 "github.com/eclipse/paho.mqtt.golang"
	"log"
)

var DevicesData = make(chan *DeviceData, 100)

type DeviceData struct {
	ID       int         `json:"uniq_id"`
	Name     string      `json:"device_name"`
	Category Category    `json:"category"`
	Data     interface{} `json:"data"`
}

type Category string

const (
	Unknown     Category = ""
	Temperature Category = "temperature"
	Humidity             = "humidity"
	Motion               = "motion"
)

var messagePubHandler mqtt2.MessageHandler = func(client mqtt2.Client, msg mqtt2.Message) {
	var data *DeviceData
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Println("can't unmarshal json", err)
	} else {
		DevicesData <- data
	}
}

var connectHandler mqtt2.OnConnectHandler = func(client mqtt2.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt2.ConnectionLostHandler = func(client mqtt2.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func New(address string, port int, clientID, username, password string) mqtt2.Client {
	opts := mqtt2.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", address, port))
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetDefaultPublishHandler(messagePubHandler)
	mqttClient := mqtt2.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return mqttClient
}

func Subscribe(topic string, client mqtt2.Client) {
	token := client.Subscribe(topic, 0, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s\n", topic)
}

func Publish(topic string, client mqtt2.Client, message string) {
	client.Publish(topic, 0, false, message)
}
