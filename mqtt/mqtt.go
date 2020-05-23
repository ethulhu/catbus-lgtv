// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

// Package mqtt wraps Paho MQTT with a few quality-of-life features.
package mqtt

import (
	"fmt"

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

func NewClientOptions() *mqtt.ClientOptions {
	return mqtt.NewClientOptions()
}
func NewClient(opts *mqtt.ClientOptions) mqtt.Client {
	return mqtt.NewClient(opts)
}

func URI(host string, port uint) string {
	return fmt.Sprintf("tcp://%v:%v", host, port)
}