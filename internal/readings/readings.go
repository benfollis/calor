package readings

import (
	"follis.net/internal/pubsub"
	"follis.net/internal/thermometers"
)

type ReadAcceptor interface {
	Accept(reading thermometers.Reading)
}

const Topic = "readings"

func InitializeChannel() *pubsub.PubSubber {
	ps := pubsub.Initialize(100) // we'll queue 100 messages to the channel if nobody is reading
	return ps
}