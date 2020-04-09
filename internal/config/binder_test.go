package config

import (
	"github.com/benfollis/calor/internal/readings/acceptors"
	"github.com/benfollis/calor/internal/thermometers"
	"testing"
)

func TestSimpleBinder_BindPort(t *testing.T) {
	cb := LoadedConfigBinder{}
	unbound := LoadedConfig{Port:9090}
	bound := cb.Bind(unbound)
	if (bound.Port != 9090) {
		t.Errorf("Should have had port %d", 9090)
	}
}

func TestSimpleBinder_BindConsoleAc(t *testing.T) {
	cb := LoadedConfigBinder{}
	unbound := LoadedConfig{ReadAcceptors: []ReadAcceptor{{
		Name:       "Bar",
		DriverType: "Console",
	}}}
	bound := cb.Bind(unbound)
	ra := bound.ReadAcceptors[0]
	if ra.Name != "Bar" {
		t.Errorf("Acceptor should have been named %s", "Bar")
	}
	switch ra.ReadAcceptor.(type) {
	case acceptors.ConsoleAcceptor:
		return
	default:
		t.Errorf("Should have been a console read acceptor")
	}
}

func TestSimpleBinder_BindZeroK(t *testing.T) {
	cb := LoadedConfigBinder{}
	unbound := LoadedConfig{Thermometers: []ThermometerConfig{{
		Name:           "Cod",
		DriverType:     "ZeroKelvin",
		UpdateInterval: 1,
		Options:        nil,
	}}}
	bound := cb.Bind(unbound)
	therm := bound.Thermometers[0]
	if therm.UpdateInterval != 1 {
		t.Errorf("Should have had update interval of %d", 1)
	}
	if therm.Name != "Cod" {
		t.Errorf("Sould have had name of %s", "Cod")
	}
	switch therm.Thermometer.(type) {
	case thermometers.ZeroKelvin:
		return
	default:
		t.Errorf("Should have been a zero kelvin thermometer")
	}
}
