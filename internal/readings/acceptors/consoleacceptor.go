package acceptors

import (
	"fmt"
	"follis.net/internal/thermometers"
)

type ConsoleAcceptor struct {
	MyName string
}

func (crs ConsoleAcceptor) Name() string {
	return crs.MyName
}

func (crs ConsoleAcceptor) Accept(reading thermometers.Reading) {
	fmt.Println( reading.Name,"at time", reading.Time, "of", reading.Temp, reading.Unit)
}
