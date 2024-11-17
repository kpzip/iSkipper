package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Token  string `json:"token"`
	UserId string `json:"user_id"`
}

func getConfig() (*Config, error) {
	data, err := os.ReadFile("./config.json")
	if err != nil {
		return nil, err
	}
	stringData := string(data)
	var cfg *Config = nil
	err = json.Unmarshal([]byte(stringData), &cfg)
	return cfg, err
}
