package config

import (
	"github.com/benfollis/calor/internal/database"
	"github.com/benfollis/calor/internal/readings"
	"github.com/benfollis/calor/internal/thermometers"
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
