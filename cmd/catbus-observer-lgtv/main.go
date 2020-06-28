// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

package main

// TODO: subscribe to App and Volume, and if they're set to invalid values, set them to the real ones.

import (
	"context"
	"path"
	"sort"
	"strconv"
	"strings"
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

			publishAppNames(config, client)
		},
		DisconnectHandler: func(client catbus.Client, err error) {
			log := logger.Background()
			log.AddField("broker-uri", config.BrokerURI)
			if err != nil {
				log.AddField("error", err)
			}
			log.Info("disconnected from Catbus")
		},
	})

	go func() {
		log := logger.Background()
		log.AddField("broker-uri", config.BrokerURI)
		log.Info("connecting to Catbus")
		if err := client.Connect(); err != nil {
			log.WithError(err).Fatal("could not connect to Catbus")
		}
	}()

	for {
		log, ctx := logger.FromContext(context.Background())
		ctx, _ = context.WithTimeout(ctx, 10*time.Second)

		log.AddField("tv", config.TV.Host)

		log.Info("connecting to TV")
		tv, err := lgtv.Dial(ctx, config.TV.Host, lgtv.DefaultOptions)
		if err != nil {
			log.WithError(err).Info("could not connect to TV")
			continue
		}
		log.Info("connected to TV")

		if _, err := tv.Register(ctx, config.TV.Key); err != nil {
			log.WithError(err).Error("could not register with TV")
			tv.Close()
			continue
		}
		log.Info("registered with TV")

		tv.SubscribeApp(func(app lgtv.App) {
			log, _ := log.Fork(context.Background())
			log.AddField("app-id", app.ID)

			if app.ID == "" {
				log.Info("empty app ID, TV about to turn off")
				return
			}

			name, ok := config.AppNameForID(app.ID)
			if !ok {
				name = app.ID
			}
			log.AddField("app-name", name)

			log.AddField("topic", config.Topics.App)
			if err := client.Publish(config.Topics.App, catbus.Retain, name); err != nil {
				log.WithError(err).Error("could not publish to Catbus")
				return
			}
			log.Info("published to Catbus")
		})

		tv.SubscribeVolume(func(v lgtv.Volume) {
			log, _ := log.Fork(context.Background())
			log.AddField("volume", v.Percent)
			log.AddField("topic", config.Topics.Volume)

			if err := client.Publish(config.Topics.Volume, catbus.Retain, strconv.Itoa(v.Percent)); err != nil {
				log.WithError(err).Error("could not publish to Catbus")
				return
			}
			log.Info("published to Catbus")
		})

		log.Info("waiting for TV to disconnect")
		if err := tv.Wait(); err != nil {
			log.WithError(err).Error("disconnected from TV")
		} else {
			log.Info("disconnected from TV")
		}
	}
}

func publishAppNames(config *config.Config, client catbus.Client) {
	log := logger.Background()

	var appNames []string
	for appName := range config.Apps {
		appNames = append(appNames, appName)
	}
	sort.Strings(appNames)
	appNamesTopic := path.Join(config.Topics.App, "values")
	if err := client.Publish(appNamesTopic, catbus.Retain, strings.Join(appNames, "\n")); err != nil {
		log.WithError(err).Error("could not publish app names to Catbus")
		return
	}
	log.Info("published app names to Catbus")
}
