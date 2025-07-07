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
	mux.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Printf("Starting go Server at port: %v", port)
	log.Fatal(srvr.ListenAndServe())
}
