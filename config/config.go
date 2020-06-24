// SPDX-FileCopyrightText: 2020 Ethel Morgan
//
// SPDX-License-Identifier: MIT

// Package config holds the definition & utilities for config.json.
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type (
	Config struct {
		BrokerURI string `json:"mqttBroker"`

		TV struct {
			Host string `json:"host"`
			Key  string `json:"key"`
		} `json:"tv"`

		Topics struct {
			App       string `json:"app"`
			AppValues string `json:"appValues"`
			Power     string `json:"power"`
			Volume    string `json:"volume"`
		} `json:"topics"`

		Apps map[string]string `json:"apps"`
	}
)

func (c *Config) AppNameForID(id string) (string, bool) {
	for name, id2 := range c.Apps {
		if id2 == id {
			return name, true
		}
	}
	return "", false
}

func Load(path string) (*Config, error) {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	config := &Config{}
	if err := json.Unmarshal(src, config); err != nil {
		return nil, fmt.Errorf("could not parse JSON: %w", err)
	}
	return config, nil
}
