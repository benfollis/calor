package web

import (
	"fmt"
	"follis.net/internal/config"
	"follis.net/internal/database"
	"follis.net/internal/utils"
	"net/http"
	"strconv"
)

type WebConfig struct {
	DB database.CalorDB
}

func Init(config config.BoundConfig) {
	webConf := WebConfig{
		DB: config.Database,
	}
	latestHandler := LatestReadingGenerator(webConf)
	http.HandleFunc("/latest/", latestHandler)
	fmt.Println("Starting web server on port", config.Port)
	addr := ":" + strconv.Itoa(config.Port)
	err := http.ListenAndServe(addr, nil)
	utils.Check(err)
}