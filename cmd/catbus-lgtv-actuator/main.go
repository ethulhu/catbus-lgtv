// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"strconv"
	"time"

	"go.eth.moe/catbus"
	"go.eth.moe/catbus-lgtv/config"
	"go.eth.moe/catbus-lgtv/lgtv"
	"go.eth.moe/flag"
	"go.eth.moe/logger"
)

var (
	configPath = flag.Custom("config-path", "", "path to config.json", flag.RequiredString)
)

func main() {
	flag.Parse()

	configPath := (*configPath).(string)

	log, _ := logger.FromContext(context.Background())

	config, err := config.Load(configPath)
	if err != nil {
		log.AddField("config-path", configPath)
		log.WithError(err).Fatal("could not load config")
	}

	client := catbus.NewClient(config.BrokerURI, catbus.ClientOptions{
		ConnectHandler: func(client catbus.Client) {
			log := logger.Background()
			log.AddField("broker-uri", config.BrokerURI)
			log.Info("connected to Catbus")

			if err := client.Subscribe(config.Topics.App, setApp(config)); err != nil {
				log := log.WithError(err)
				log.AddField("topic", config.Topics.App)
				log.Error("could not subscribe to topic")
			}
			if err := client.Subscribe(config.Topics.Volume, setVolume(config)); err != nil {
				log := log.WithError(err)
				log.AddField("topic", config.Topics.Volume)
				log.Error("could not subscribe to topic")
			}
			if err := client.Subscribe(config.Topics.Power, setPower(config)); err != nil {
				log := log.WithError(err)
				log.AddField("topic", config.Topics.Power)
				log.Error("could not subscribe to topic")
			}
			log.Info("subscribed to topics")
		},
		DisconnectHandler: func(client catbus.Client, err error) {
			log := logger.Background()
			log.AddField("broker-uri", config.BrokerURI)
			if err != nil {
				log.AddField("error", err)
			}
			log.Warning("disconnected from Catbus")
		},
	})

	log.AddField("broker-uri", config.BrokerURI)
	log.Info("connecting to Catbus")
	if err := client.Connect(); err != nil {
		log.WithError(err).Fatal("could not connect to Catbus")
	}
}

func setApp(config *config.Config) catbus.MessageHandler {
	return func(_ catbus.Client, msg catbus.Message) {
		log, ctx := logger.FromContext(context.Background())

		log.AddField("app-name", msg.Payload)

		appID, ok := config.Apps[msg.Payload]
		if !ok {
			log.Warning("got invalid app name")
			return
		}

		log.AddField("app-id", appID)

		ctx, _ = context.WithTimeout(ctx, 5*time.Second)
		tv, err := lgtv.Dial(ctx, config.TV.Host, lgtv.DefaultOptions)
		if err != nil {
			log.WithError(err).Warning("could not connect to TV")
			return
		}
		if _, err := tv.Register(ctx, config.TV.Key); err != nil {
			log.WithError(err).Error("could not register with TV")
			tv.Close()
			return
		}

		if err := tv.SetApp(ctx, appID); err != nil {
			log.WithError(err).Error("could not set app")
			return
		}
		log.Info("set app")
	}
}
func setVolume(config *config.Config) catbus.MessageHandler {
	return func(_ catbus.Client, msg catbus.Message) {
		log, ctx := logger.FromContext(context.Background())

		log.AddField("volume-raw", msg.Payload)

		volume, err := strconv.Atoi(msg.Payload)
		if err != nil {
			log.WithError(err).Warning("got invalid volume")
			return
		}
		if volume < 0 {
			volume = 0
		}
		if volume > 100 {
			volume = 100
		}

		log.AddField("volume", volume)

		ctx, _ = context.WithTimeout(ctx, 5*time.Second)
		tv, err := lgtv.Dial(ctx, config.TV.Host, lgtv.DefaultOptions)
		if err != nil {
			log.WithError(err).Warning("could not connect to TV")
			return
		}
		if _, err := tv.Register(ctx, config.TV.Key); err != nil {
			log.WithError(err).Error("could not register with TV")
			tv.Close()
			return
		}

		if err := tv.SetVolume(ctx, volume); err != nil {
			log.WithError(err).Error("could not set volume")
			return
		}
		log.Info("set volume")
	}
}
func setPower(config *config.Config) catbus.MessageHandler {
	return func(_ catbus.Client, msg catbus.Message) {
		log, ctx := logger.FromContext(context.Background())

		log.AddField("power", msg.Payload)

		if msg.Payload != "off" {
			log.Info("power is not \"off\", ignoring")
			return
		}

		ctx, _ = context.WithTimeout(ctx, 5*time.Second)
		tv, err := lgtv.Dial(ctx, config.TV.Host, lgtv.DefaultOptions)
		if err != nil {
			log.WithError(err).Warning("could not connect to TV")
			return
		}
		if _, err := tv.Register(ctx, config.TV.Key); err != nil {
			log.WithError(err).Error("could not register with TV")
			tv.Close()
			return
		}

		if err := tv.TurnOff(ctx); err != nil {
			log.WithError(err).Error("could not turn TV off")
			return
		}
		log.Info("turned TV off")
	}
}
