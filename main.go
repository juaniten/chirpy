package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"
	const rootPath = "."
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(rootPath)))))
	serveMux.HandleFunc("GET /metrics", apiCfg.hitsHandler)
	serveMux.HandleFunc("POST /reset", apiCfg.resetHitsHandler)
	serveMux.HandleFunc("GET /healthz", healthHandler)

	log.Printf("Listing on port %s: serving files from `%s`.\n", port, rootPath)
	server := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
	server.ListenAndServe()
}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	bodyContents := "OK"
	w.Write([]byte(bodyContents))
}

func (cfg *apiConfig) hitsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())))
}
func (cfg *apiConfig) resetHitsHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("Hits resseted"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		cfg.fileserverHits.Add(1)
	})
}
