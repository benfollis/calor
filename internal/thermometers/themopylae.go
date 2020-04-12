package thermometers

import (
	"encoding/json"
	"github.com/benfollis/calor/internal/utils"
	"io/ioutil"
	"net/http"
	"time"
)

/*
A Thermometer that reads from a Thermopylae
https://github.com/benfollis/thermopylae
remote thermometer
 */

type jsonReading struct {
	Unit string `json:"unit""`
	Temperature float64 `json:"temperature"`
}

type Thermopylae struct {
	Name string
	Url string
}

func (leonidas Thermopylae) Read() (Reading, error) {
	resp, err := http.Get(leonidas.Url)
	if err != nil {
		return Reading{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	utils.CheckLog(err)
	if err != nil {
		return Reading{}, err
	}
	var jr jsonReading
	parseErr := json.Unmarshal(body, &jr)
	utils.CheckLog(parseErr)
	if parseErr != nil {
		return Reading{}, parseErr
	}
	return Reading{
		Temp: jr.Temperature,
		Unit: jr.Unit,
		Name: leonidas.Name,
		Time: time.Now(),
	}, nil
}
