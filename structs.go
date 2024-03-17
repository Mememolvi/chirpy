package main

import "net/http"

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
}

type chirp struct {
	Body string `json:"body"`
}

type errorResponse struct {
	Error string `json:"error"`
}
type responseValid struct {
	Valid bool `json:"valid"`
}
