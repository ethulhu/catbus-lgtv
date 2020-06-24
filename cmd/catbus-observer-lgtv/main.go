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
		ConnectHandler: func(client *catbus.Client) {
			log.Printf("connected to Catbus %v", config.BrokerURI)

		},
		DisconnectHandler: func(client *catbus.Client, err error) {
			log.Printf("disconnected from Catbus %s: %v", config.BrokerURI, err)
		},
	})

	go func() {
		log.Printf("connecting to Catbus %v", config.BrokerURI)
		if err := client.Connect(); err != nil {
			log.Fatalf("could not connect to Catbus: %v", err)
		}
	}()

	for {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		tv, err := lgtv.Dial(ctx, config.TV.Host, lgtv.DefaultOptions)
		if err != nil {
			log.Printf("could not connect to TV: %v", err)
			continue
		}
		log.Print("connected to TV")

		if _, err := tv.Register(ctx, config.TV.Key); err != nil {
			log.Printf("could not register with TV: %v", err)
			tv.Close()
			continue
		}
		log.Print("registered with TV")

		tv.SubscribeApp(func(app lgtv.App) {
			if app.ID == "" {
				// This means it's about to turn off.
				return
			}

			name, ok := config.AppNameForID(app.ID)
			if !ok {
				name = app.ID
			}

			if err := client.Publish(config.Topics.App, catbus.Retain, name); err != nil {
				log.Printf("could not publish to %q: %v", config.Topics.App, err)
			}
		})

		tv.SubscribeVolume(func(v lgtv.Volume) {
			if err := client.Publish(config.Topics.Volume, catbus.Retain, strconv.Itoa(v.Percent)); err != nil {
				log.Printf("could not publish to %q: %v", config.Topics.Volume, err)
			}
		})

		if err := tv.Wait(); err != nil {
			log.Printf("disconnected from TV: %v", err)
		}
	}
}
