package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/ethulhu/mqtt-lgtv-bridge/config"
	"github.com/ethulhu/mqtt-lgtv-bridge/lgtv"
	"github.com/ethulhu/mqtt-lgtv-bridge/mqtt"
	"github.com/ethulhu/mqtt-lgtv-bridge/wol"
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

	log.Printf("connecting to MQTT broker %v:%v", cfg.BrokerHost, cfg.BrokerPort)
	broker := mqtt.NewClient(cfg.BrokerHost, cfg.BrokerPort)

	tv := lgtv.NewClient(cfg.TV.Host, lgtv.DefaultOptions)

	tv.SetAppHandler(func(app lgtv.App) {
		name, ok := cfg.AppNameForID(app.ID)
		if !ok {
			name = app.ID
		}
		broker.Publish(cfg.TopicInput, mqtt.AtLeastOnce, mqtt.Retain, name)
	})
	tv.SetVolumeHandler(func(volume lgtv.Volume) {
		broker.Publish(cfg.TopicVolume, mqtt.AtLeastOnce, mqtt.Retain, volume.Level)
	})

	broker.Subscribe(cfg.TopicVolume, mqtt.AtLeastOnce, handler(func(payload string) error {
		volume, err := strconv.Atoi(payload)
		if err != nil || !(0 <= volume && volume <= 100) {
			// Try to set the volume on the bus to the actual value, but don't worry if it fails.
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if v, err := tv.Volume(ctx); err == nil {
				broker.Publish(cfg.TopicVolume, mqtt.AtLeastOnce, mqtt.Retain, strconv.Itoa(v.Level))
			}
			return fmt.Errorf("volume must be an integer within [0,100], found %v", payload)
		}

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		return tv.SetVolume(ctx, volume)
	}))
	broker.Subscribe(cfg.TopicInput, mqtt.AtLeastOnce, handler(func(payload string) error {
		appID, ok := cfg.Apps[payload]
		if !ok {
			// Try to set the app on the bus to the actual value, but don't worry if it fails.
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if a, err := tv.App(ctx); err == nil {
				broker.Publish(cfg.TopicVolume, mqtt.AtLeastOnce, mqtt.Retain, a.Name)
			}
			return fmt.Errorf("invalid channel %v", payload)
		}

		// TODO: return the old app / actual app if it was invalid so the bus is valid.
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		return tv.SetApp(ctx, lgtv.App{ID: appID})
	}))

	broker.Subscribe(cfg.TopicPower, mqtt.AtLeastOnce, handler(func(payload string) error {
		if payload == "on" && !tv.IsConnected() {
			log.Print("turning TV on")
			mac, _ := net.ParseMAC(cfg.TV.MAC)
			if err := wol.Wake(mac); err != nil {
				return fmt.Errorf("could not turn on TV: %w", err)
			}
		}

		if payload == "off" && tv.IsConnected() {
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			return tv.TurnOff(ctx)
		}

		return nil
	}))

	// loop forever.
	key := cfg.TV.Key
	for {
		log.Print("connecting to TV")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		err := tv.Connect(ctx)
		for err != nil {
			ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
			err = tv.Connect(ctx)
			log.Printf("could not connect to TV: %v", err)
		}
		if key, err = tv.Register(ctx, key); err != nil {
			log.Printf("could not register with TV: %v", err)
		}
		log.Print("connected to TV")

		if err := tv.Err(); err != nil {
			log.Printf("disconnected from TV: %v", err)
		}
		broker.Publish(cfg.TopicPower, mqtt.AtLeastOnce, mqtt.Retain, "off")
	}
}

func handler(f func(string) error) mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload())

		if err := f(payload); err != nil && !errors.Is(err, lgtv.ErrNotConnected) {
			log.Print(err)
		}
	}

}
