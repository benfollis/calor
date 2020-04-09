package thermometers

import "time"

// A Reading denotes a temperature, and the name and time
// the reading was taken
type Reading struct {
	Temp float64
	Unit string
	Name string
	Time time.Time
}

// All Thermometers must be able to take a reading
type Thermometer interface {
	Read() Reading
}
