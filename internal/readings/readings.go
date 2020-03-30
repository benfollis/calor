package readings

import (
	"follis.net/internal/thermometers"
	"github.com/cskr/pubsub"
)

type ReadAcceptor interface {
	Accept(reading thermometers.Reading)
}

const Topic = "readings"

func InitializeChannel() *pubsub.PubSub {
	ps := pubsub.New(100) // we'll queue 100 messages to the channel if nobody is reading
	return ps
}