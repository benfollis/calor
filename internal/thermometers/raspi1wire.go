package thermometers

import (
	"calor/internal/utils"
	"github.com/yryz/ds18b20"
	"time"
)


type Raspi1Wire struct{
	Name string
	SensorId string
}

func (rw1 Raspi1Wire) Read() Reading{
	temp, err := ds18b20.Temperature(rw1.SensorId)
	utils.Check(err)
	reading := Reading{
		Temp: temp,
		Unit: "C",
		Name: rw1.Name,
		Time: time.Now(),
	}
	return reading
}