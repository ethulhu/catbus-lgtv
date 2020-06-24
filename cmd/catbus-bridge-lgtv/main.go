// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.eth.moe/catbus-lgtv/config"
	"go.eth.moe/catbus-lgtv/lgtv"
	"go.eth.moe/catbus"
)

var (
	configPath = flag.String("config-path", "", "path to config.json")
)

func main() {
	flag.Parse()

	if *configPath == "" {
		log.Fatal("must set --config-path")
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	tv := lgtv.NewClient(cfg.TV.Host, lgtv.DefaultOptions)

	brokerURI := fmt.Sprintf("tcp://%v:%v", cfg.BrokerHost, cfg.BrokerPort)
	client := catbus.NewClient(brokerURI, catbus.ClientOptions{
		ConnectHandler: func(client *catbus.Client) {
			log.Printf("connected to Catbus %v", brokerURI)

			if err := client.Subscribe(cfg.TopicApp, setApp(cfg, tv)); err != nil {
				log.Printf("could not subscribe to %q: %v", cfg.TopicApp, err)
			}
			if err := client.Subscribe(cfg.TopicPower, setPower(cfg, tv)); err != nil {
				log.Printf("could not subscribe to %q: %v", cfg.TopicPower, err)
			}
			if err := client.Subscribe(cfg.TopicVolume, setVolume(tv)); err != nil {
				log.Printf("could not subscribe to %q: %v", cfg.TopicVolume, err)
			}

			var appNames []string
			for name := range cfg.Apps {
				appNames = append(appNames, name)
			}
			sort.Strings(appNames)
			if err := client.Publish(cfg.TopicAppValues, catbus.Retain, strings.Join(appNames, "\n")); err != nil {
				log.Printf("could not publish to %q: %v", cfg.TopicAppValues, err)
			}
		},
		DisconnectHandler: func(client *catbus.Client, err error) {
			log.Printf("disconnected from Catbus %s: %v", brokerURI, err)
		},
	})

	go func() {
		log.Printf("connecting to Catbus %v", brokerURI)
		if err := client.Connect(); err != nil {
			log.Fatalf("could not connect to Catbus: %v", err)
		}
	}()

	tv.SetAppHandler(func(app lgtv.App) {
		name, ok := cfg.AppNameForID(app.ID)
		if !ok {
			name = app.ID
		}
		if err := client.Publish(cfg.TopicApp, catbus.Retain, name); err != nil {
			log.Printf("could not publish to %q: %v", cfg.TopicApp, err)
		}
	})
	tv.SetVolumeHandler(func(volume lgtv.Volume) {
		if err := client.Publish(cfg.TopicVolume, catbus.Retain, fmt.Sprintf("%v", volume.Level)); err != nil {
			log.Printf("could not publish to %q: %v", cfg.TopicVolume, err)
		}
	})

	// loop forever.
	key := cfg.TV.Key
	for {
		log.Print("connecting to TV")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		err := tv.Connect(ctx)
		for err != nil {
			log.Printf("could not connect to TV: %v", err)
			ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
			err = tv.Connect(ctx)
		}
		if key, err = tv.Register(ctx, key); err != nil {
			log.Printf("could not register with TV: %v", err)
		}
		log.Print("connected to TV")

		if err := tv.SubscribeApp(); err != nil {
			log.Printf("could not subscribe to app events: %v", err)
		}
		if err := tv.SubscribeVolume(); err != nil {
			log.Printf("could not subscribe to volume events: %v", err)
		}

		if err := tv.Err(); err != nil {
			log.Printf("disconnected from TV: %v", err)
		}
	}
}

func setVolume(tv lgtv.Client) catbus.MessageHandler {
	return func(_ *catbus.Client, msg catbus.Message) {
		volume, err := strconv.Atoi(string(msg.Payload()))
		if err != nil {
			return
		}
		if volume < 0 {
			volume = 0
		}
		if volume > 100 {
			volume = 100
		}

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		if err := tv.SetVolume(ctx, volume); err != nil && !errors.Is(err, lgtv.ErrNotConnected) {
			log.Print(err)
		}
	}
}
func setApp(cfg *config.Config, tv lgtv.Client) catbus.MessageHandler {
	return func(_ *catbus.Client, msg catbus.Message) {
		appID, ok := cfg.Apps[string(msg.Payload())]
		if !ok {
			return
		}

		// TODO: return the old app / actual app if it was invalid so the bus is valid.
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		if err := tv.SetApp(ctx, lgtv.App{ID: appID}); err != nil && !errors.Is(err, lgtv.ErrNotConnected) {
			log.Print(err)
		}
	}
}
func setPower(cfg *config.Config, tv lgtv.Client) catbus.MessageHandler {
	return func(_ *catbus.Client, msg catbus.Message) {
		switch string(msg.Payload()) {
		case "off":
			if tv.IsConnected() {
				ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
				if err := tv.TurnOff(ctx); err != nil && !errors.Is(err, lgtv.ErrNotConnected) {
					log.Print(err)
				}
			}
		default:
			return
		}
	}
}
