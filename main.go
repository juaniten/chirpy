package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/juaniten/chirpy/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error opening postgres database.")
	}

	const port = "8080"
	const rootPath = "."
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
	}

	serveMux := http.NewServeMux()
	fileserverHandler := http.FileServer(http.Dir(rootPath))
	serveMux.Handle("/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(fileserverHandler)))
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	serveMux.HandleFunc("GET /api/healthz", healthHandler)
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	serveMux.HandleFunc("POST /api/users", apiCfg.handlerUsers)

	log.Printf("Listing on port %s: serving files from `%s`.\n", port, rootPath)
	server := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
	server.ListenAndServe()
}
