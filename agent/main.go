package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type Metric struct {
	Hostname  string  `json:"hostname"`
	CPU       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	Timestamp int64   `json:"timestamp"`
}

func collectMetrics() Metric {
	hostname, _ := os.Hostname()
	cpuPercent, _ := cpu.Percent(time.Second, false)
	vmStat, _ := mem.VirtualMemory()

	return Metric{
		Hostname:  hostname,
		CPU:       cpuPercent[0],
		Memory:    vmStat.UsedPercent,
		Timestamp: time.Now().Unix(),
	}
}

func sendMetrics(m Metric, backendURL string) {
	data, _ := json.Marshal(m)

	resp, err := http.Post(backendURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("❌ Error sending metrics:", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("✅ Sent metrics: Host=%s CPU=%.2f%% Mem=%.2f%%", m.Hostname, m.CPU, m.Memory)
}

func main() {
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		log.Fatal("❌ BACKEND_URL environment variable is not set")
	}

	for {
		metric := collectMetrics()
		sendMetrics(metric, backendURL)
		time.Sleep(5 * time.Second)
	}
}
