package main

import (
	"fmt"
	"log"
	"net/http"
)

type newServer http.ServeMux


func main(){
	const port = "8080"

	mux := http.NewServeMux()
	
	srvr := &http.Server{
		Handler: mux,
		Addr: ":" + port,
	}
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", handler)
	fmt.Printf("Starting go Server at port: %v", port)
	log.Fatal(srvr.ListenAndServe())
}

func handler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
