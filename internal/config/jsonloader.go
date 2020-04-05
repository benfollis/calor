package config

import (
	"encoding/json"
	"follis.net/internal/utils"
	"io/ioutil"
)

type JsonLoader struct{}

func (loader JsonLoader) Load(filepath string) LoadedConfig {
	contents, err := ioutil.ReadFile(filepath)
	utils.Check(err)
	var config LoadedConfig
	parseErr := json.Unmarshal(contents, &config)
	utils.Check(parseErr)
	return config
}
