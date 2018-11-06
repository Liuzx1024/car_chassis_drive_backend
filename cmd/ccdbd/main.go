package main

import (
	"flag"
	"fmt"
)

func main() {
	configFile := flag.String("config", "config.json", "Load config from json file")
	flag.Parse()
	fmt.Println(*configFile)
	if err := loadConfig(*configFile); err != nil {
		panic(err)
	}
	if err := buildDeviceTree(); err != nil {
		panic(err)
	}
	if err := runBackend(); err != nil {
		panic(err)
	}
	if err := runAPIServer(); err != nil {
		panic(err)
	}
}

//TODO IMPLEMENT
func buildDeviceTree() error {
	return nil
}

//TODO IMPLEMENT
func runAPIServer() error {
	return nil
}

//TODO IMPLEMENT
func runBackend() error {
	return nil
}
