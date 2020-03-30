package thermometers

import "time"

type Reading struct {
	Temp float64
	Unit string
	Name string
	Time time.Time
}

type Thermometer interface {
	Read() Reading
}
