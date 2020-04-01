package config

import "follis.net/internal/thermometers"

type BoundThermometer struct {
	Name string
	Thermometer thermometers.Thermometer
	UpdateInterval int
}

type BoundConfig struct {
	Thermometers []BoundThermometer
	Database DatabaseConfig
	Port int
}

type ConfigBinder interface {
	Bind(config LoadedConfig) BoundConfig
}
