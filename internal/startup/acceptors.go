package startup

import (
	"fmt"
	"github.com/benfollis/calor/internal/config"
	"github.com/benfollis/calor/internal/pubsub"
	"github.com/benfollis/calor/internal/readings"
	"github.com/benfollis/calor/internal/thermometers"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)


func spawnAcceptor(ch chan interface{}, reader readings.ReadAcceptor, group * sync.WaitGroup) {
	defer group.Done()
	for message := range ch {
		reading := message.(thermometers.Reading)
		reader.Accept(reading)
	}
}

func StartAcceptors(config config.BoundConfig, ps *pubsub.PubSub) sync.WaitGroup {
	fmt.Println("Starting acceptors")
	var rwg sync.WaitGroup
	// if any of the read acceptors is sqlite, make sure the db exists
	for _, bra := range config.ReadAcceptors {
		acceptor := bra.ReadAcceptor
		ch := ps.Subscribe(readings.Topic)
		rwg.Add(1)
		go spawnAcceptor(ch, acceptor, &rwg)
		// now start up the acceptor
	}
	return rwg
}