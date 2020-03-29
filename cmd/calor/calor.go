package main

import (
	"flag"
	"fmt"
	"follis.net/internal/config"
)

func main() {
	pathFlag := flag.String("config", "/etc/calor/calor.json", "The location of calor's json config file")
	flag.Parse()
	var configLoader config.ConfigLoader
	configLoader = config.JsonLoader{}
	loadedConfig := configLoader.Load(*pathFlag)
	var configBinder config.ConfigBinder
	configBinder = config.SimpleBinder{}
	bound := configBinder.Bind(loadedConfig)
	for _, boundTherm := range bound.Thermometers {
		thermometer := boundTherm.Thermometer
		reading := thermometer.Read()
		fmt.Println("Reading is ", reading.Temp)
	}
}
