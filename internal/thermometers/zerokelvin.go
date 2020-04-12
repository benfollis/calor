/**
	A themometer that always reads 0 degrees kelvin
 */
package thermometers

import (
	"time"
)

// A ZeroKelvin thermometer is a Dummy Thermometer that
// always reports 0K
type ZeroKelvin struct {
	Name string
}

func (zk ZeroKelvin) Read() (Reading, error) {
	reading := Reading{
		Temp: 0,
		Unit: "K",
		Name: zk.Name,
		Time: time.Now(),
	}
	return reading, nil
}
