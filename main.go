package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct{
	fileServerHits atomic.Int32
}

func main(){
	const port = "8080"

	mux := http.NewServeMux()

	apicfg := &apiConfig{
		fileServerHits: atomic.Int32{},
	}
	
	srvr := &http.Server{
		Handler: mux,
		Addr: ":" + port,
	}

	mux.Handle("/app/", apicfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/healthz", handler)
	mux.HandleFunc("/metrics", apicfg.writeHandler)
	mux.HandleFunc("/reset", apicfg.resetHandler)
	fmt.Printf("Starting go Server at port: %v", port)
	log.Fatal(srvr.ListenAndServe())
}

func handler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) writeHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	value := "Hits: " + fmt.Sprintf("%v", cfg.fileServerHits.Load())
	w.Write([]byte(value))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset"))
}
