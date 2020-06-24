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
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.eth.moe/catbus-lgtv/config"
	"go.eth.moe/catbus-lgtv/lgtv"
	"go.eth.moe/catbus-lgtv/mqtt"
	"go.eth.moe/catbus-lgtv/wol"
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
	brokerOptions := mqtt.NewClientOptions()
	brokerOptions.AddBroker(brokerURI)
	brokerOptions.SetAutoReconnect(true)
	brokerOptions.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		log.Printf("disconnected from MQTT broker %s: %v", brokerURI, err)
	})
	brokerOptions.SetOnConnectHandler(func(broker mqtt.Client) {
		log.Printf("connected to MQTT broker %v", brokerURI)

		_ = broker.Subscribe(cfg.TopicApp, mqtt.AtLeastOnce, setApp(cfg, tv))
		_ = broker.Subscribe(cfg.TopicPower, mqtt.AtLeastOnce, setPower(cfg, tv))
		_ = broker.Subscribe(cfg.TopicVolume, mqtt.AtLeastOnce, setVolume(tv))
	})

	log.Printf("connecting to MQTT broker %v", brokerURI)
	broker := mqtt.NewClient(brokerOptions)
	_ = broker.Connect()

	var appNames []string
	for name := range cfg.Apps {
		appNames = append(appNames, name)
	}
	sort.Strings(appNames)
	_ = broker.Publish(cfg.TopicAppValues, mqtt.AtLeastOnce, mqtt.Retain, strings.Join(appNames, "\n"))

	tv.SetAppHandler(func(app lgtv.App) {
		name, ok := cfg.AppNameForID(app.ID)
		if !ok {
			name = app.ID
		}
		_ = broker.Publish(cfg.TopicApp, mqtt.AtLeastOnce, mqtt.Retain, name)
	})
	tv.SetVolumeHandler(func(volume lgtv.Volume) {
		_ = broker.Publish(cfg.TopicVolume, mqtt.AtLeastOnce, mqtt.Retain, fmt.Sprintf("%v", volume.Level))
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
		_ = broker.Publish(cfg.TopicPower, mqtt.AtLeastOnce, mqtt.Retain, "on")

		if err := tv.SubscribeApp(); err != nil {
			log.Printf("could not subscribe to app events: %v", err)
		}
		if err := tv.SubscribeVolume(); err != nil {
			log.Printf("could not subscribe to volume events: %v", err)
		}

		if err := tv.Err(); err != nil {
			log.Printf("disconnected from TV: %v", err)
		}
		_ = broker.Publish(cfg.TopicPower, mqtt.AtLeastOnce, mqtt.Retain, "off")
	}
}

func setVolume(tv lgtv.Client) mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
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
func setApp(cfg *config.Config, tv lgtv.Client) mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
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
func setPower(cfg *config.Config, tv lgtv.Client) mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		switch string(msg.Payload()) {
		case "on":
			mac, _ := net.ParseMAC(cfg.TV.MAC)
			if err := wol.Wake(mac); err != nil {
				log.Printf("could not turn on TV: %v", err)
			}
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
