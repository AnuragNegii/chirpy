package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct{
	ID uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Email string `json:"email"`
}

func (apiConfig *apiConfig) handleUser(w http.ResponseWriter, r *http.Request){
	type returnVals struct{
		Email string `json:"email"`
	}

	var rV returnVals

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rV)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("cant decode the request: %v", err), nil)
		return
	}
	
	user, err := apiConfig.db.CreateUser(r.Context(), rV.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("can create user: %v", err), nil)
		return
	}

	respondWithJson(w, http.StatusCreated, User{
		ID: user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email: user.Email,
	})
}
