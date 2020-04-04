package startup

import (
	"database/sql"
	"fmt"
	"follis.net/internal/config"
	"follis.net/internal/database"
	"follis.net/internal/pubsub"
	"follis.net/internal/readings"
	"follis.net/internal/readings/acceptors"
	"follis.net/internal/thermometers"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)


func createDB(sqla acceptors.SqLiteAcceptor) {
	// only have SQLITE with files for now
	//create the tables
	db, err := sql.Open("sqlite3", sqla.DBFile)
	if err != nil {
		panic(err)
	}
	database.CreateTable(db)
	db.Close()
}

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
		switch acceptor.(type) {
		case acceptors.SqLiteAcceptor:
			sqla := acceptor.(acceptors.SqLiteAcceptor)
			createDB(sqla)
		default:
		}
		ch := ps.Subscribe(readings.Topic)
		rwg.Add(1)
		go spawnAcceptor(ch, acceptor, &rwg)
		// now start up the acceptor
	}
	return rwg
}