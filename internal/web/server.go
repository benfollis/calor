package web

import (
	"fmt"
	"github.com/benfollis/calor/internal/config"
	"github.com/benfollis/calor/internal/database"
	"github.com/benfollis/calor/internal/utils"
	"net/http"
	"strconv"
)

type WebConfig struct {
	DB database.CalorDB
	Thermometers []config.BoundThermometer
}

func Init(config config.BoundConfig) {
	webConf := WebConfig{
		DB: config.Database,
		Thermometers: config.Thermometers,
	}
	latest := calorGenerator{
		config: webConf,
		handler:latest,
	}
	between := calorGenerator{
		config:webConf,
		handler:between,
	}
	discovery := calorGenerator{
		config:webConf,
		handler:discover,
	}
	http.Handle("/latest/", latest)
	http.Handle("/between/", between)
	http.Handle("/discovery", discovery)
	fmt.Println("Starting web server on port", config.Port)
	addr := ":" + strconv.Itoa(config.Port)
	err := http.ListenAndServe(addr, nil)
	utils.CheckLog(err)
}