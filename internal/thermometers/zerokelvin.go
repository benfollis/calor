/**
	A themometer that always reads 0 degrees kelvin
 */
package thermometers

import (
	"time"
)

type ZeroKelvin struct {
	Name string
}

func (zk ZeroKelvin) Read() Reading {
	reading := Reading{
		Temp: 0,
		Unit: "K",
		Name: zk.Name,
		Time: time.Now(),
	}
	return reading;
}
