package config

import (
	"encoding/json"
	"github.com/benfollis/calor/internal/utils"
	"io/ioutil"
)

// A config loader that reads from a JSON file
type JsonLoader struct{}

func (loader JsonLoader) Load(filepath string) LoadedConfig {
	contents, err := ioutil.ReadFile(filepath)
	utils.CheckPanic(err)
	var config LoadedConfig
	parseErr := json.Unmarshal(contents, &config)
	utils.CheckLog(parseErr)
	return config
}
