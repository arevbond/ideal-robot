package mqtt

import (
	"fmt"
	mqtt2 "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt2.MessageHandler = func(client mqtt2.Client, msg mqtt2.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
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
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt2.NewClient(opts)
	return client
}

func Subscribe(topic string, client mqtt2.Client) {
	token := client.Subscribe(topic, 0, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s\n", topic)
}

func Publish(topic string, client mqtt2.Client, message string) {
	client.Publish(topic, 0, false, message)
}
