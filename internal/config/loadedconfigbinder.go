package config

import (
	"github.com/benfollis/calor/internal/database"
	"github.com/benfollis/calor/internal/readings/acceptors"
	"github.com/benfollis/calor/internal/thermometers"
)

// A LoadedConfig bnder accepts a LoadedConfig
// and produces a Bound Config from it.
type LoadedConfigBinder struct{}

func (sb LoadedConfigBinder) Bind(config LoadedConfig) BoundConfig {
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
	switch config.Database.DriverType {
	case "Sqlite":
		boundDB = database.CreateSqliteDB(config.Database.File)
	case "Postgres":
		boundDB = database.CreatePostgresDB(config.Database.Username, config.Database.Password, config.Database.Host)
	}
	boundAcceptors := make([]BoundReadAcceptor, len(config.ReadAcceptors))
	for index, unboundRA := range config.ReadAcceptors {
		bound := BoundReadAcceptor{
			Name: unboundRA.Name,
		}
		switch unboundRA.DriverType {
		case "Console":
			bound.ReadAcceptor = acceptors.ConsoleAcceptor{MyName: unboundRA.Name}
		case "DB":
			bound.ReadAcceptor = acceptors.DBAcceptor{
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
