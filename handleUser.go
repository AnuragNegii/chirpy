package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AnuragNegii/chirpy/internal/auth"
	"github.com/AnuragNegii/chirpy/internal/database"

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
		Password string `json:"password"`
	}
	var rV returnVals
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rV)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("cant decode the request: %v", err), nil)
		return
	}
	newString, err := auth.HashPassword(rV.Password)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
		return
	}
	rV.Password = newString 
	user, err := apiConfig.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: rV.Email,
		HashedPassword: rV.Password,
	})
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
