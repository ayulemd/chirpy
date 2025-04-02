package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)

		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileserverHits.Load()
	fmt.Fprintf(w, "Hits: %d", hits)
}

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	hits := cfg.fileserverHits.Load()
	fmt.Fprintf(w, "Hits: %d", hits)
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	cfg := apiConfig{}
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fileServer := http.FileServer(http.Dir(filepathRoot))

	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(fileServer)))
	mux.Handle("/app/assets/", http.StripPrefix("/app", cfg.middlewareMetricsInc(fileServer)))
	mux.HandleFunc("GET /api/healthz", handlerHealthCheck)
	mux.HandleFunc("GET /api/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /api/reset", cfg.handlerResetMetrics)

	server.ListenAndServe()
}

func handlerHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
