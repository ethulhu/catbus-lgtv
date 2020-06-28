// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"flag"
	"log"
	"strconv"
	"time"

	"go.eth.moe/catbus"
	"go.eth.moe/catbus-lgtv/config"
	"go.eth.moe/catbus-lgtv/lgtv"
)

var (
	configPath = flag.String("config-path", "", "path to config.json")
)

func main() {
	flag.Parse()

	if *configPath == "" {
		log.Fatal("must set --config-path")
	}

	config, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	client := catbus.NewClient(config.BrokerURI, catbus.ClientOptions{
		ConnectHandler: func(client catbus.Client) {
			log.Printf("connected to Catbus %v", config.BrokerURI)

			if err := client.Subscribe(config.Topics.App, setApp(config)); err != nil {
				log.Printf("could not subscribe to %q: %v", config.Topics.App, err)
			}
			if err := client.Subscribe(config.Topics.Volume, setVolume(config)); err != nil {
				log.Printf("could not subscribe to %q: %v", config.Topics.Volume, err)
			}
			if err := client.Subscribe(config.Topics.Power, setPower(config)); err != nil {
				log.Printf("could not subscribe to %q: %v", config.Topics.Power, err)
			}
		},
		DisconnectHandler: func(client catbus.Client, err error) {
			log.Printf("disconnected from Catbus %s: %v", config.BrokerURI, err)
		},
	})

	log.Printf("connecting to Catbus %v", config.BrokerURI)
	if err := client.Connect(); err != nil {
		log.Fatalf("could not connect to Catbus: %v", err)
	}
}

func setApp(config *config.Config) catbus.MessageHandler {
	return func(_ catbus.Client, msg catbus.Message) {
		appID, ok := config.Apps[msg.Payload]
		if !ok {
			log.Printf("got invalid app %q", msg.Payload)
			return
		}

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		tv, err := lgtv.Dial(ctx, config.TV.Host, lgtv.DefaultOptions)
		if err != nil {
			log.Printf("could not connect to TV: %v", err)
			return
		}
		if _, err := tv.Register(ctx, config.TV.Key); err != nil {
			log.Printf("could not register with TV: %v", err)
			tv.Close()
			return
		}

		if err := tv.SetApp(ctx, appID); err != nil {
			log.Printf("could not set app to %q: %v", appID, err)
		}
	}
}
func setVolume(config *config.Config) catbus.MessageHandler {
	return func(_ catbus.Client, msg catbus.Message) {
		volume, err := strconv.Atoi(msg.Payload)
		if err != nil {
			log.Printf("got invalid volume %q", volume)
			return
		}
		if volume < 0 {
			volume = 0
		}
		if volume > 100 {
			volume = 100
		}

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		tv, err := lgtv.Dial(ctx, config.TV.Host, lgtv.DefaultOptions)
		if err != nil {
			log.Printf("could not connect to TV: %v", err)
			return
		}
		if _, err := tv.Register(ctx, config.TV.Key); err != nil {
			log.Printf("could not register with TV: %v", err)
			tv.Close()
			return
		}

		if err := tv.SetVolume(ctx, volume); err != nil {
			log.Printf("could not set volume to %v: %v", volume, err)
		}
	}
}
func setPower(config *config.Config) catbus.MessageHandler {
	return func(_ catbus.Client, msg catbus.Message) {
		if msg.Payload != "off" {
			return
		}

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		tv, err := lgtv.Dial(ctx, config.TV.Host, lgtv.DefaultOptions)
		if err != nil {
			log.Printf("could not connect to TV: %v", err)
			return
		}
		if _, err := tv.Register(ctx, config.TV.Key); err != nil {
			log.Printf("could not register with TV: %v", err)
			tv.Close()
			return
		}

		if err := tv.TurnOff(ctx); err != nil {
			log.Printf("could not turn TV off: %v", err)
		}
	}
}
