package config

import (
	"calor/internal/database"
	"calor/internal/readings"
	"calor/internal/thermometers"
)

type BoundThermometer struct {
	Name string
	Thermometer thermometers.Thermometer
	UpdateInterval int
}

type BoundReadAcceptor struct {
	Name string
	ReadAcceptor readings.ReadAcceptor
}

type BoundConfig struct {
	Thermometers []BoundThermometer
	Database database.CalorDB
	Port int
	ReadAcceptors []BoundReadAcceptor
}

type ConfigBinder interface {
	Bind(config LoadedConfig) BoundConfig
}
