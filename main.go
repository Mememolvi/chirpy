package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	ac := apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()

	mux.Handle("/app/*", ac.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /api/metrics", func(w http.ResponseWriter, r *http.Request) {
		h := fmt.Sprintf("Hits: %v", ac.fileserverHits)
		w.Write([]byte(h))
	})
	mux.HandleFunc("GET /api/reset", ac.reset)
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		s := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", ac.fileserverHits)
		w.Write([]byte(s))
	})
	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		chirp := chirp{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&chirp)
		if err != nil {
			er := errorResponse{
				Error: "Something went wrong",
			}
			w.WriteHeader(500)
			b, _ := json.Marshal(er)
			w.Write(b)
		}

		if len(chirp.Body) > 140 {
			er := errorResponse{
				Error: "Chirp is too long",
			}
			w.WriteHeader(400)
			b, _ := json.Marshal(er)
			w.Write(b)
		}

		rv := responseValid{
			Valid: true,
		}
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(rv)
		w.Write(b)
	})

	corsMux := middlewareCors(mux)
	http.ListenAndServe(":8080", corsMux)
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
