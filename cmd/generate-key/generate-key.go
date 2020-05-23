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

	"github.com/ethulhu/catbus-lgtv/config"
	"github.com/ethulhu/catbus-lgtv/lgtv"
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

	tv := lgtv.NewClient(cfg.TV.Host, lgtv.DefaultOptions)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err = tv.Connect(ctx); err != nil {
		log.Fatalf("could not connect to TV: %v", err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
	key, err := tv.Register(ctx, "")
	if err != nil {
		log.Fatalf("could not register to TV: %v", err)
	}
	fmt.Printf("key: %v\n", key)
}