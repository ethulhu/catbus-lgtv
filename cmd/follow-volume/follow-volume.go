// Binary follow-volume listens for and prints changes to the volume on a WebOS LG TV.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ethulhu/mqtt-lgtv-bridge/config"
	"github.com/ethulhu/mqtt-lgtv-bridge/lgtv"
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

	tv.SetVolumeHandler(func(v lgtv.Volume) {
		fmt.Println(v)
	})

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	if err := tv.Connect(ctx); err != nil {
		log.Fatalf("could not connect to TV: %v", err)
	}
	if _, err := tv.Register(ctx, cfg.TV.Key); err != nil {
		log.Fatalf("could not register with TV: %v", err)
	}

	_ = tv.SubscribeVolume()

	// block forever.
	select {}
}
