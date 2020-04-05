package web

import (
	"encoding/json"
	"fmt"
	"follis.net/internal/thermometers"
	"net/http"
	"strings"
)

func LatestReadingGenerator(config WebConfig) func(w http.ResponseWriter, r *http.Request) {
	db := config.DB
	handler := func(w http.ResponseWriter, r *http.Request) {
		url := r.URL
		path := url.Path
		chunks := strings.Split(path, "/")
		// path should be of the form /latest/<thermometer name>
		if len(chunks) != 2 {
			w.WriteHeader(400)
		}
		name := chunks[1]
		reading := db.Latest(name)
		empty := thermometers.Reading{}
		if reading == empty {
			w.WriteHeader(404)
		}
		encoded, _ := json.Marshal(reading)
		fmt.Fprint(w, encoded)
	}
	return handler
}
