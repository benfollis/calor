package config

import "follis.net/internal/thermometers"

type BoundThermometer struct {
	Name string
	Thermometer thermometers.Thermometer
	UpdateInterval int
}

type BoundConfig struct {
	Thermometers []BoundThermometer
	Port int
}

type ConfigBinder interface {
	Bind(config LoadedConfig) BoundConfig
}
