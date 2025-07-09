package main

import (
	"encoding/json"
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
	mux.HandleFunc("GET /api/healthz", handler)
	mux.HandleFunc("GET /admin/metrics", apicfg.writeHandler)
	mux.HandleFunc("POST /admin/reset", apicfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", hadnle_Chirp)
	fmt.Printf("Starting go Server at port: %v", port)
	log.Fatal(srvr.ListenAndServe())
}

func handler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) writeHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
	`, cfg.fileServerHits.Load())))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset"))
}

func hadnle_Chirp(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body) 
	var params parameters
	err := decoder.Decode(&params)
	if err != nil{
		log.Printf("error decoding parameters: %v", err)
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusInternalServerError)
		return
	}
	
	if len(params.Body) > 140{
		resp := map[string]string{"error": "Chirp is too long"}
		w.Header().Set("Content-Type", "application/`json:")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := map[string]bool{"valid":true}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
