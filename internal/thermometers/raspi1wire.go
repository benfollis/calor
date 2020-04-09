package thermometers

import (
	"github.com/benfollis/calor/internal/utils"
	"github.com/yryz/ds18b20"
	"time"
)

// A Raspi1Wire is a thermometer that takes it's readings from a
// Linux kernel 1 wire therm interface thermometer
type Raspi1Wire struct{
	Name string
	SensorId string
}

func (rw1 Raspi1Wire) Read() Reading{
	temp, err := ds18b20.Temperature(rw1.SensorId)
	utils.CheckLog(err)
	reading := Reading{
		Temp: temp,
		Unit: "C",
		Name: rw1.Name,
		Time: time.Now(),
	}
	return reading
}