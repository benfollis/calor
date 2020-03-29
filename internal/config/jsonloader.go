package config

import (
	"encoding/json"
	"io/ioutil"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type JsonLoader struct {}


func (loader JsonLoader) Load(filepath string) LoadedConfig {
	contents, err := ioutil.ReadFile(filepath)
	check(err)
	var config LoadedConfig
	parseErr := json.Unmarshal(contents, &config)
	check(parseErr)
	return config
}
