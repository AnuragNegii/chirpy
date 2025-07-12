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
	type returnVals struct{
		Valid bool `json:"valid"`
	}
	decoder := json.NewDecoder(r.Body) 
	var params parameters
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength{
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJson(w, http.StatusOK, returnVals{
		Valid:true,
	})
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error){
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	type errorResponse struct{
		Error string `json:"error"`
	}

	respondWithJson(w, code, errorResponse{
		Error:msg,
	})

}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}){

	w.Header().Set("Content-Type","application/json")
	data, err := json.Marshal(payload)
	if err != nil{
		log.Printf("Error marshalling json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
