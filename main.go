package main

import (
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
	fileserver := http.FileServer(http.Dir(rootPath))
	serveMux.Handle("/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(fileserver)))
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetHitsHandler)
	serveMux.HandleFunc("GET /api/healthz", healthHandler)

	log.Printf("Listing on port %s: serving files from `%s`.\n", port, rootPath)
	server := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
	server.ListenAndServe()
}
func (cfg *apiConfig) resetHitsHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("Hits resseted"))
}
