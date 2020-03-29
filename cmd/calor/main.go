package main

import (
	"fmt"
	"follis.net/internal/thermometers"
)

func main() {
	var zeroK thermometers.Thermometer
	zeroK = thermometers.ZeroKelvin{}
	reading := zeroK.Read()
	fmt.Println("Reading is ", reading.Temp)
}
