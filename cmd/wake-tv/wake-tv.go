// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

// Binary wake-tv wakes a given WebOS LG TV.
package main

import (
	"flag"
	"log"
	"net"

	"github.com/ethulhu/catbus-lgtv/config"
	"github.com/ethulhu/catbus-lgtv/wol"
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

	mac, err := net.ParseMAC(cfg.TV.MAC)
	if err != nil {
		log.Fatalf("invalid MAC address %q: %v", cfg.TV.MAC, err)
	}
	if err := wol.Wake(mac); err != nil {
		log.Fatalf("could not wake TV: %v", err)
	}
}