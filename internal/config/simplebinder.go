package config

import (
	"follis.net/internal/readings"
	"follis.net/internal/thermometers"
)

type SimpleBinder struct{}

func (sb SimpleBinder) Bind(config LoadedConfig) BoundConfig {
	numTherms := len(config.Thermometers)
	boundTherms := make([]BoundThermometer, numTherms)
	for index, unboundTherm := range config.Thermometers {
		bound := BoundThermometer{
			Name:           unboundTherm.Name,
			UpdateInterval: unboundTherm.UpdateInterval,
		}
		switch unboundTherm.DriverType {
		case "ZeroKelvin":
			bound.Thermometer = thermometers.ZeroKelvin{Name:unboundTherm.Name}
		}
		boundTherms[index] = bound
	}
	for _, unboundRA := range config.ReadAcceptors {
		bound := BoundReadAcceptor{
			Name: unboundRA.Name,
		}
		switch unboundRA.DriverType {
		case "Console":
			bound.ReadAcceptor = readings.ConsoleAcceptor{MyName: unboundRA.Name}
		case "Sqlite":
			bound.ReadAcceptor = readings.SqLiteAcceptor{
			MyName:unboundRA.Name,
			DBFile: config.Database.File,
			}
		}
	}
	return BoundConfig{
		Thermometers: boundTherms,
		Port:         config.Port,
		Database:	  config.Database,
	}
}
