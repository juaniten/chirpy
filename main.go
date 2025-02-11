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
	platform       string
	jwtSecret      string
	polkaKey       string
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
		platform:       os.Getenv("PLATFORM"),
		jwtSecret:      os.Getenv("JWT_SECRET"),
		polkaKey:       os.Getenv("POLKA_KEY"),
	}

	serveMux := http.NewServeMux()
	fileserverHandler := http.FileServer(http.Dir(rootPath))
	serveMux.Handle("/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(fileserverHandler)))

	serveMux.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	serveMux.HandleFunc("GET /api/healthz", healthHandler)
	serveMux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	serveMux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)

	serveMux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsGet)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpGet)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)

	serveMux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	serveMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeToken)

	serveMux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhook)

	log.Printf("Listing on port %s: serving files from `%s`.\n", port, rootPath)
	server := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
	server.ListenAndServe()
}
