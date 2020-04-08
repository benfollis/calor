package web

import (
	"fmt"
	"calor/internal/config"
	"calor/internal/database"
	"calor/internal/utils"
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
	latestHandler := LatestGenerator(webConf)
	betweenHandler := BetweenGenerator(webConf)
	discoveryHandler := DiscoveryGenerator(webConf)
	http.HandleFunc("/latest/", latestHandler)
	http.HandleFunc("/between/", betweenHandler)
	http.HandleFunc("/discovery", discoveryHandler)
	fmt.Println("Starting web server on port", config.Port)
	addr := ":" + strconv.Itoa(config.Port)
	err := http.ListenAndServe(addr, nil)
	utils.Check(err)
}