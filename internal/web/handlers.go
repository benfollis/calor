package web

import (
	"encoding/json"
	"fmt"
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
		if len(chunks) != 3 {
			w.WriteHeader(400)
		}
		name := chunks[2]
		fmt.Println("Received request for thermometer", name)
		reading, err := db.Latest(name)
		if err != nil {
			w.WriteHeader(404)
			fmt.Fprint(w, "Not Found")
			return
		}
		encoded, _ := json.Marshal(reading)
		stringEncoded := string(encoded)
		fmt.Println(stringEncoded)
		fmt.Fprint(w, stringEncoded)
	}
	return handler
}
