// Binary list-apps lists the apps for a given WebOS LG TV.
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

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	log.Printf("connecting to TV")
	if err := tv.Connect(ctx); err != nil {
		log.Fatalf("could not connect to TV: %v", err)
	}
	if _, err := tv.Register(ctx, cfg.TV.Key); err != nil {
		log.Fatalf("could not register with TV: %v", err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	apps, err := tv.ListApps(ctx)
	if err != nil {
		log.Fatalf("could not get apps: %v", err)
	}
	for _, app := range apps {
		fmt.Printf("%v: %v\n", app.ID, app.Name)
	}
}
