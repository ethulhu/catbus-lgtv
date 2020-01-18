// Binary list-apps lists the apps for a given WebOS LG TV.
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ethulhu/mqtt-lgtv-bridge/config"
	"github.com/mjanser/lgtv"
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
		log.Fatalf("failed to load config from %v: %v", *configPath, err)
	}

	url := fmt.Sprintf("ws://%v:3000", cfg.TV.Host)
	tv := lgtv.NewDefaultClient(url, cfg.TV.Key)

	log.Printf("connecting to TV %v", url)
	if err := tv.Connect(); err != nil {
		log.Fatalf("failed to connect to TV: %v", err)
	}

	apps, err := tv.GetApps()
	if err != nil {
		log.Fatalf("failed to get apps: %v", err)
	}

	for _, app := range apps {
		fmt.Printf("%v: %v\n", app.ID, app.Title)
	}
}
