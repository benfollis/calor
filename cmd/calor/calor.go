package main

import (
	"flag"
	"follis.net/internal/config"
	"follis.net/internal/readings"
	"follis.net/internal/thermometers"
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
	ps := readings.InitializeChannel()

	readers := make([]readings.ReadAcceptor, 1)
	readers[0] = readings.ConsoleAcceptor{}

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
					ps.Pub(reading, readings.Topic)
				}
			}
		}(&twg)
		reading := thermometer.Read()
		ps.Pub(reading, readings.Topic)
	}
	// fire up our reader wait group
	var rwg sync.WaitGroup
	// for every reader start a go routine subscribing to the readers
	for _, reader := range readers {
		rwg.Add(1)
		go func(group *sync.WaitGroup) {
			defer group.Done();
			ch := ps.Sub(readings.Topic)
			for message := range ch {
				reading := message.(thermometers.Reading)
				reader.Accept(reading)
			}
		}(&rwg)
	}



	twg.Wait()

}
