package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const rootPath = "."

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(rootPath))))
	serveMux.HandleFunc("/healthz", healthHandler)

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
