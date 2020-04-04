package main

import (
	"flag"
	"follis.net/internal/config"
	"follis.net/internal/startup"
	"follis.net/internal/pubsub"
	"follis.net/internal/readings"
	"sync"
	"time"
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
	// start up the pub sub channels
	ps := pubsub.Initialize(100)
	rwg := startup.StartAcceptors(bound, ps)
	//start up the producers
	var twg sync.WaitGroup
	// for testing, lets publish a message. We'll move this to a ticker soon
	for _, boundTherm := range bound.Thermometers {
		thermometer := boundTherm.Thermometer
		ticker := time.NewTicker(time.Duration(boundTherm.UpdateInterval) * time.Second)
		twg.Add(1)
		go func (group *sync.WaitGroup) {
			defer group.Done()
			for {
				select {
				case <-ticker.C:
					reading := thermometer.Read()
					ps.Publish(readings.Topic, reading)
				}
			}
		}(&twg)
	}
	rwg.Wait()
	twg.Wait()

}
