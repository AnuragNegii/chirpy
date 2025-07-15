package main

import (
	"context"
	"fmt"
	"net/http"
)


func (apiConfig *apiConfig) ResetHandle(w http.ResponseWriter, r *http.Request){
	if apiConfig.platform != "dev"{
		respondWithError(w, http.StatusForbidden, "you are not allowed to reset", nil)
		return
	}
	apiConfig.fileServerHits.Store(0)
	err := apiConfig.db	.ResetUsers(context.Background())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("rest tabel users: %v", err), nil)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset"))
}
