// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

// Binary set-volume sets the volume on a WebOS LG TV.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	"go.eth.moe/catbus-lgtv/config"
	"go.eth.moe/catbus-lgtv/lgtv"
)

var (
	volumePercent = flag.Int("volume-percent", -1, "volume percent to set")
	configPath    = flag.String("config-path", "", "path to config.json")
)

func main() {
	flag.Parse()

	if *configPath == "" || *volumePercent == -1 {
		log.Fatal("must set --config-path and --volume-percent")
	}
	if !(0 <= *volumePercent && *volumePercent <= 100) {
		log.Fatalf("--volume-percent must be within 0 and 100, got %v", *volumePercent)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	log.Print("connecting to TV")
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	tv, err := lgtv.Dial(ctx, cfg.TV.Host, lgtv.DefaultOptions)
	if err != nil {
		log.Fatalf("could not dial TV: %v", err)
	}

	log.Print("registering with TV")
	if _, err := tv.Register(ctx, cfg.TV.Key); err != nil {
		log.Fatalf("could not register with TV: %v", err)
	}

	log.Print("setting volume")
	tv.SetVolume(ctx, *volumePercent)
}
