package readings

import (
	"fmt"
	"follis.net/internal/thermometers"
)

type ConsoleAcceptor struct {
}

func (crs ConsoleAcceptor) Accept(reading thermometers.Reading) {
	fmt.Println( reading.Name,"at time", reading.Time, "of", reading.Temp, reading.Unit)
}
