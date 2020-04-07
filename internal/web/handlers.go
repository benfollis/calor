package web

import (
	"encoding/json"
	"fmt"
	"follis.net/internal/database"
	"net/http"
	"strconv"
	"strings"
)

func respondWithData(data interface{}, w http.ResponseWriter) {
	encoded, _ := json.Marshal(data)
	stringEncoded := string(encoded)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(encoded)))
	fmt.Fprint(w, stringEncoded)
}

func respondNotFound(w http.ResponseWriter) {
	w.WriteHeader(404)
	fmt.Fprint(w, "Not Found")
}

func DiscoveryGenerator(config WebConfig) func(w http.ResponseWriter, r *http.Request) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		names := make([]string, len(config.Thermometers))
		for index, therm := range config.Thermometers {
			names[index] = therm.Name
		}
		respondWithData(names, w);
	}
	return handler
}

func LatestGenerator(config WebConfig) func(w http.ResponseWriter, r *http.Request) {
	db := config.DB
	handler := func(w http.ResponseWriter, r *http.Request) {
		url := r.URL
		path := url.Path
		chunks := strings.Split(path, "/")
		// path should be of the form /latest/<thermometer name>
		if len(chunks) != 3 {
			w.WriteHeader(400)
		}
		name := chunks[2]
		fmt.Println("Received request for thermometer", name)
		reading, err := db.Latest(name)
		if err != nil {
			respondNotFound(w)
			return
		}
		respondWithData(reading, w);
	}
	return handler
}

func BetweenGenerator(config WebConfig) func(w http.ResponseWriter, r *http.Request) {
	db := config.DB
	handler := func(w http.ResponseWriter, r *http.Request) {
		url := r.URL
		path := url.Path
		chunks := strings.Split(path, "/")
		// path should be of the form /between/<thermometer name>?start=<unixTime>&end=<unixTime>
		if len(chunks) != 3 {
			w.WriteHeader(400)
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
		if err != nil {
			respondNotFound(w)
			return
		}
		respondWithData(readings, w);
	}
	return handler
}
