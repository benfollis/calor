package main

import (
	"database/sql"
	"flag"
	"fmt"
	"follis.net/internal/config"
	"follis.net/internal/database"
	"follis.net/internal/readings"
	"follis.net/internal/thermometers"
	_ "github.com/mattn/go-sqlite3"
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
	// TODO, factor this stuff into methods
	readers := make([]readings.ReadAcceptor, 1)
	readers[0] = readings.ConsoleAcceptor{}
	if bound.Database.File != "" {
		// only have SQLITE with files for now
		//create the tables
		db, err := sql.Open("sqlite3", bound.Database.File)
		if err != nil {
			panic(err)
		}
		database.CreateTable(db)
		db.Close()
		dbReader := readings.SqLiteAcceptor{DBFile:bound.Database.File}
		readers = append(readers, dbReader)
	}
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
	// fire up our reader wait group
	var rwg sync.WaitGroup
	// for every reader start a go routine subscribing to the readers
	for _, reader := range readers {
		rwg.Add(1)
		go func(group *sync.WaitGroup) {
			defer group.Done();
			ch := ps.Subscribe(readings.Topic)
			for message := range ch {
				fmt.Println("Got message", ch)
				reading := message.(thermometers.Reading)
				reader.Accept(reading)
			}
		}(&rwg)
	}



	twg.Wait()

}
