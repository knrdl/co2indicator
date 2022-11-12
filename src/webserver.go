package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func startWebserver(bind string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		latestMeasurement.Lock()
		j, _ := json.Marshal(&latestMeasurement.value)
		latestMeasurement.Unlock()
		w.Write(j)
	})
	log.Printf("Webserver is running on %v\n", bind)
	log.Fatal(http.ListenAndServe(bind, nil))
}
