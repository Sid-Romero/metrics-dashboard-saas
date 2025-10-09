package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Metric struct {
	Hostname  string  `json:"hostname"`
	CPU       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	Timestamp int64   `json:"timestamp"`
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var m Metric
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received metric: Host=%s CPU=%.2f%% Mem=%.2f%% Time=%s",
		m.Hostname, m.CPU, m.Memory, time.Unix(m.Timestamp, 0))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func main() {
	http.HandleFunc("/metrics", metricsHandler)

	port := ":8080"
	log.Printf("Metrics API running on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
