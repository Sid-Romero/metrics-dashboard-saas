package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	dbpkg "metrics-dashboard/backend/internal/db"
)

type Metric struct {
	Hostname  string  `json:"hostname"`
	CPU       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	Timestamp int64   `json:"timestamp"`
}

type Server struct {
	Store *dbpkg.Store
}

func (s *Server) MetricsPOST(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST", http.StatusMethodNotAllowed)
		return
	}

	var m Metric
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ts := time.Unix(m.Timestamp, 0).UTC()

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := s.Store.InsertMetric(ctx, m.Hostname, m.CPU, m.Memory, ts); err != nil {
		http.Error(w, "db insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) HealthGET(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	if err := s.Store.Pool.Ping(ctx); err != nil {
		http.Error(w, "db not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// JSON test 
func (s *Server) MetricsGET(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	host := q.Get("hostname")
	limit := int32(100)
	if l := q.Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = int32(n)
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	rows, err := s.Store.Recent(ctx, host, limit)
	if err != nil {
		http.Error(w, "db query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rows)
}
