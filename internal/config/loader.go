package config

// A ConfigLoader takes a config file and produces a LoadedConfig
// Currently we just have the json loader
type ConfigLoader interface {
	Load(filepath string) LoadedConfig
}

// A ThermometerConfig is the basic config for a thermometer
// The options map exist to allow key/value pairs that are specific to a
// particular driver type
type ThermometerConfig struct {
	Name string
	DriverType string
	UpdateInterval int // time in seconds to take a reading
	Options map[string] string
}

// A DabaaseConfig contains all config needed to describe a DB connection
type DatabaseConfig struct {
	DriverType string
	File string
	Host string
	Port string
	Username string
	Password string
}

// A ReadAcceptor contains all config needed to construct a particular type of
// read acceptor
type ReadAcceptor struct {
	Name string
	DriverType string
}

// A LoadedConfig is a struct declaring what we want,
// but does not actually contain any valid calor types
type LoadedConfig struct {
	Thermometers []ThermometerConfig
	Database DatabaseConfig
	Port int // TCP port to start listening on
	ReadAcceptors []ReadAcceptor // the things readings get output to
}
