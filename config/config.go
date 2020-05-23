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
		BrokerHost string `json:"broker_host"`
		BrokerPort uint   `json:"broker_port"`

		TV struct {
			Host string `json:"host"`
			MAC  string `json:"mac"`
			Key  string `json:"key"`
		} `json:"tv"`

		TopicPower     string `json:"topic_power"`
		TopicApp       string `json:"topic_app"`
		TopicAppValues string `json:"topic_app_values"`
		TopicVolume    string `json:"topic_volume"`

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