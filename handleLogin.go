package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AnuragNegii/chirpy/internal/auth"
)



func (apiConfig *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request){
	type Params struct{
		Email string `json:"email"`	
		Password string `json:"password"`
		Expires_In_Seconds int `json:"expires_in_seconds"`
	}

	var params Params

	decoder  := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
		return
	}

	getUser, err := apiConfig.db.CheckLogin(r.Context(), params.Email)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
		return 
	}
	err = auth.CheckPasswordHash(params.Password, getUser.HashedPassword)		
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("%v", err), nil)
		return
	}

	respondWithJson(w, http.StatusOK, User{
		ID: getUser.ID,
		Created_at: getUser.CreatedAt,
		Updated_at: getUser.UpdatedAt,
		Email: getUser.Email,
	})

}
