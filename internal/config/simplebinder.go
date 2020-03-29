package config

import "follis.net/internal/thermometers"

type SimpleBinder struct{}

func (sb SimpleBinder) Bind(config LoadedConfig) BoundConfig {
	numTherms := len(config.Thermometers)
	boundTherms := make([]BoundThermometer, numTherms)
	for _, unboundTherm := range config.Thermometers {
		switch unboundTherm.DriverType {
		case "ZeroKelvin":
			bound := BoundThermometer{
				Name:        unboundTherm.Name,
				Thermometer: thermometers.ZeroKelvin{},
			}
			boundTherms = append(boundTherms, bound)
		}
	}
	return BoundConfig{
		Thermometers: boundTherms,
		Port:         config.Port,
	}
}
