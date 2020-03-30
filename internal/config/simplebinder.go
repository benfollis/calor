package config

import "follis.net/internal/thermometers"

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
	return BoundConfig{
		Thermometers: boundTherms,
		Port:         config.Port,
	}
}
