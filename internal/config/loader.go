package config


type ConfigLoader interface {
	Load(filepath string) LoadedConfig
}

type ThermometerConfig struct {
	Name string
	DriverType string
	UpdateInterval int // time in seconds to take a reading
	Options map[string] string
}

type LoadedConfig struct {
	Thermometers []ThermometerConfig
	Port int // TCP port to start listening on
}
