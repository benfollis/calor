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

type DatabaseConfig struct {
	Driver string
	File string
	Host string
	Port string
	Username string
	Password string
}

type LoadedConfig struct {
	Thermometers []ThermometerConfig
	Database DatabaseConfig
	Port int // TCP port to start listening on
}
