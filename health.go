package main

import "net/http"

func healthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	bodyContents := "OK"
	w.Write([]byte(bodyContents))
}
