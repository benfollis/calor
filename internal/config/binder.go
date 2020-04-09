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

// A BoundConfig represents a confuration that has actual
// calor structs loaded into memory
type BoundConfig struct {
	Thermometers []BoundThermometer
	Database database.CalorDB
	Port int
	ReadAcceptors []BoundReadAcceptor
}

// A ConfigBinder is a type that can take some LoadedConfig and produce a Bound Config
// From it. Currently it's a fairly useless abstraction, but eventually we'll want to be able
// bind configs from other sources, and hence we don't want to directly work with
// LoadedConfigs
type ConfigBinder interface {
	Bind(config LoadedConfig) BoundConfig
}
