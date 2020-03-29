/**
	A themometer that always reads 0 degrees kelvin
 */
package thermometers

type ZeroKelvin struct {
}

func (zk ZeroKelvin) Read() Reading {
	reading := Reading { 0.0, "K" }
	return reading;
}
