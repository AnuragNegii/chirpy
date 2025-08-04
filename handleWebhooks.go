package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AnuragNegii/chirpy/internal/auth"
	"github.com/google/uuid"
)

type Webhook struct{
	Event string `json:"event"`
	Data struct{
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (apiConfig *apiConfig) handleWebhooks(w http.ResponseWriter, r *http.Request){
	var webhook Webhook
	defer r.Body.Close()
	getKey, err := auth.GetAPIKey(r.Header)
	if err != nil{
		respondWithError(w, 401, fmt.Sprintf("%v", err), nil)
		return
	}
	if getKey != apiConfig.polkaKey{
		respondWithError(w, 401, fmt.Sprintf("%v", err), nil)
		return
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&webhook)
	if err != nil{
		respondWithError(w, 204, fmt.Sprintf("%v", err), nil)
		return
	}
	if webhook.Event != "user.upgraded" {
		respondWithError(w, 204, fmt.Sprintf("%v", err), nil)
		return
	}
	userID, err := uuid.Parse(webhook.Data.UserID)
	if err != nil {
		respondWithError(w, 204, fmt.Sprintf("%v", err), nil)
		return
	}
	err = apiConfig.db.UpgradeUsers(r.Context(), userID)
	if err != nil{
		respondWithError(w, 404, fmt.Sprintf("%v", err), nil)
		return
	}
	w.WriteHeader(204)	
}
