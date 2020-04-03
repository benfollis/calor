package readings

import (
	"follis.net/internal/thermometers"
)

type ReadAcceptor interface {
	Name() string
	Accept(reading thermometers.Reading)
}

const Topic = "readings"
