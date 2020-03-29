package thermometers

type Reading struct {
	Temp float64
	Unit string
}

type Thermometer interface {
	Read() Reading
}
