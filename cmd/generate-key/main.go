// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

// Binary generate-key queries a given WebOS LG TV for a new connection key.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

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

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("could not load config from %v: %v", *configPath, err)
	}

	log.Printf("connecting to TV")
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	tv, err := lgtv.Dial(ctx, cfg.TV.Host, lgtv.DefaultOptions)
	if err != nil {
		log.Fatalf("could not dial TV: %v", err)
	}

	log.Printf("registering with TV")
	key, err := tv.Register(ctx, "")
	if err != nil {
		log.Fatalf("could not register to TV: %v", err)
	}
	fmt.Printf("key: %v\n", key)
}
