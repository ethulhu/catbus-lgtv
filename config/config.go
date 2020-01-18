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

		TopicPower  string `json:"topic_power"`
		TopicInput  string `json:"topic_input"`
		TopicVolume string `json:"topic_volume"`

		Apps map[string]string `json:"apps"`
	}
)

func Load(path string) (*Config, error) {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	config := &Config{}
	if err := json.Unmarshal(src, config); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return config, nil
}
