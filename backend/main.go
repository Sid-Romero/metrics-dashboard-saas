package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"metrics-dashboard/backend/handlers"
	dbpkg "metrics-dashboard/backend/internal/db"
)

func main() {
	ctx := context.Background()

	store, err := dbpkg.New(ctx)
	if err != nil {
		log.Fatalf("db init error: %v", err)
	}
	defer store.Close()

	srv := &handlers.Server{Store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			srv.MetricsPOST(w, r)
			return
		}
		srv.MetricsGET(w, r)
	})
	mux.HandleFunc("/health", srv.HealthGET)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	httpSrv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("backend listening on :%s", port)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(ctxShutdown)
}
