package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AnuragNegii/chirpy/internal/auth"
	"github.com/AnuragNegii/chirpy/internal/database"
)


func (apiConfig *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request){
	type Params struct{
		Email string `json:"email"`	
		Password string `json:"password"`
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
	
	jwtToken, err := auth.MakeJWT(getUser.ID, apiConfig.secret)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil{
	respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}	
	_, err = apiConfig.db.NewRefreshToken(r.Context(), database.NewRefreshTokenParams{
		Token: refreshToken,
		UserID: getUser.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})
	if err != nil{
	respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), nil)
	}

	respondWithJson(w, http.StatusOK, User{
		ID: getUser.ID,
		Created_at: getUser.CreatedAt,
		Updated_at: getUser.UpdatedAt,
		Email: getUser.Email,
		Token: jwtToken,
		Refresh_token: refreshToken,
	})
}
