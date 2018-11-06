package main

import (
	"encoding/json"
	"os"
)

type ccdbdConfig struct {
}

var config *ccdbdConfig

func loadConfig(path string) error {
	config = new(ccdbdConfig)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decorder := json.NewDecoder(file)
	if err := decorder.Decode(config); err != nil {
		return err
	}
	return nil
}
