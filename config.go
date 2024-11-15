package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Username string
	Password string
}

func get_config() (*Config, error) {
	data, err := os.ReadFile("./config.json")
	if err != nil {
		return nil, err
	}
	stringData := string(data)
	var cfg *Config = nil
	err = json.Unmarshal([]byte(stringData), &cfg)
	return cfg, err
}
