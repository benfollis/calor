package readings

import (
	"follis.net/internal/thermometers"
)

type ReadAcceptor interface {
	Accept(reading thermometers.Reading)
}

const Topic = "readings"
