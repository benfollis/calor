package readings

import (
	"github.com/benfollis/calor/internal/thermometers"
)

type ReadAcceptor interface {
	Name() string
	Accept(reading thermometers.Reading)
}

const Topic = "readings"
