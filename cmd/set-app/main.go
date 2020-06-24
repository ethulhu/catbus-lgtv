// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

// Binary set-app sets the running app on a WebOS LG TV.
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
	appID      = flag.String("app-id", "", "app ID to set")
	configPath = flag.String("config-path", "", "path to config.json")
)

func main() {
	flag.Parse()

	if *configPath == "" || *appID == "" {
		log.Fatal("must set --config-path and --app-id")
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

	log.Print("setting app")
	tv.SetApp(ctx, *appID)
}
