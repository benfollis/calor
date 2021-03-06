package main

import (
	"flag"
	"github.com/benfollis/calor/internal/config"
	"github.com/benfollis/calor/internal/pubsub"
	"github.com/benfollis/calor/internal/startup"
	"github.com/benfollis/calor/internal/web"
)

func main() {
	pathFlag := flag.String("config", "/etc/calor/calor.json", "The location of calor's json config file")
	flag.Parse()
	var configLoader config.ConfigLoader
	configLoader = config.JsonLoader{}
	loadedConfig := configLoader.Load(*pathFlag)
	var configBinder config.ConfigBinder
	configBinder = config.LoadedConfigBinder{}
	bound := configBinder.Bind(loadedConfig)
	bound.Database.Init() // create any DB tables needed
	// fire up our web server
	// start up the pub sub channels
	ps := pubsub.Initialize(100)
	// start our acceptors
	rwg := startup.StartAcceptors(bound, ps)
	//start up the producers
	twg := startup.StartThermometers(bound, ps)
	web.Init(bound)
	rwg.Wait()
	twg.Wait()

}
