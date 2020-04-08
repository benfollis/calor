package config

import (
	"calor/internal/database"
	"calor/internal/readings/acceptors"
	"calor/internal/thermometers"
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
			bound.Thermometer = thermometers.ZeroKelvin{Name: unboundTherm.Name}
		case "Raspi1Wire":
			options := unboundTherm.Options
			sensorId := options["SensorId"]
			bound.Thermometer = thermometers.Raspi1Wire{
				Name:     unboundTherm.Name,
				SensorId: sensorId,
			}
		}
		boundTherms[index] = bound
	}
	var boundDB database.CalorDB
	if config.Database.DriverType == "Sqlite" {
		boundDB = database.CreateSqliteDB(config.Database.File)
	}
	boundAcceptors := make([]BoundReadAcceptor, len(config.ReadAcceptors))
	for index, unboundRA := range config.ReadAcceptors {
		bound := BoundReadAcceptor{
			Name: unboundRA.Name,
		}
		switch unboundRA.DriverType {
		case "Console":
			bound.ReadAcceptor = acceptors.ConsoleAcceptor{MyName: unboundRA.Name}
		case "Sqlite":
			bound.ReadAcceptor = acceptors.SqLiteAcceptor{
				MyName: unboundRA.Name,
				DB:     boundDB,
			}
		}
		boundAcceptors[index] = bound
	}
	boundConfig := BoundConfig{
		Thermometers:  boundTherms,
		Port:          config.Port,
		ReadAcceptors: boundAcceptors,
		Database:      boundDB,
	}
	return boundConfig

}
