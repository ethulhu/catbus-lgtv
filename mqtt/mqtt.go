// Package mqtt wraps Paho MQTT with a few quality-of-life features.
package mqtt

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type (
	Client         = mqtt.Client
	Message        = mqtt.Message
	MessageHandler = mqtt.MessageHandler
)

const (
	AtMostOnce byte = iota
	AtLeastOnce
	ExactlyOnce
)

const (
	Retain = true
)

// NewClient returns a new Paho MQTT Client for a given host & port
func NewClient(host string, port uint) Client {
	brokerURI := fmt.Sprintf("tcp://%v:%v", host, port)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURI)
	opts.SetAutoReconnect(true)
	opts.SetOnConnectHandler(func(_ mqtt.Client) {
		log.Printf("connected to MQTT broker %s", brokerURI)
	})
	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		log.Printf("disconnected from MQTT broker %s: %v", brokerURI, err)
	})

	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	return client
}
