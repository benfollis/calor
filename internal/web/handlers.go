package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/benfollis/calor/internal/database"
	"net/http"
	"strconv"
	"strings"
)

// A handy function to convert data to json and respond with it
func respondWithData(data interface{}, w http.ResponseWriter) {
	encoded, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(encoded)))
	// make cors happy on the browsers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(encoded)
}

// A convenient function to respond with 404
func respondNotFound(w http.ResponseWriter) {
	w.WriteHeader(404)
	fmt.Fprint(w, "Not Found")
}


type calorGenerator struct {
	config WebConfig
	handler calorHandler
}

type calorHandler func(config WebConfig, w http.ResponseWriter, r *http.Request) (interface{}, error)

func (cg calorGenerator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if data ,err := cg.handler(cg.config, w, r); err != nil {
		respondNotFound(w)
	}else{
		respondWithData(data, w)
	}
}

func discover(config WebConfig, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	names := make([]string, len(config.Thermometers))
	for index, therm := range config.Thermometers {
		names[index] = therm.Name
	}
	return names, nil
}


func latest(config WebConfig, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	url := r.URL
	path := url.Path
	chunks := strings.Split(path, "/")
	// path should be of the form /latest/<thermometer name>
	if len(chunks) != 3 {
		return nil, errors.New("no thermometer")
	}
	db := config.DB
	name := chunks[2]
	fmt.Println("Received request for thermometer", name)
	reading, err := db.Latest(name)
	return reading, err
}

func between(config WebConfig, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	db := config.DB
	url := r.URL
	path := url.Path
	chunks := strings.Split(path, "/")
	// path should be of the form /between/<thermometer name>?start=<unixTime>&end=<unixTime>
	if len(chunks) != 3 {
		return nil, errors.New("no thermometer")
	}
	name := chunks[2]
	fmt.Println("Received between request for thermometer", name)
	query := url.Query()
	start, _ := strconv.ParseInt(query.Get("start"), 10, 64)
	end, _ := strconv.ParseInt(query.Get("end"), 10, 64)
	timestampRange := database.UnixTimestampRange{
		Begin: start,
		End:   end,
	}
	readings, err := db.Between(name, timestampRange)
	return readings, err
}

