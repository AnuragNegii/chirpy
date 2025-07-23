package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AnuragNegii/chirpy/internal/auth"
)


func (apiConfig *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request){
	type refreshToken struct{
		Token string `json:"token"`
	}
	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, 401, fmt.Sprintf("%v", err), nil)
		return
	}
		
	user, err := apiConfig.db.GetUserFromRefreshToken(r.Context(), authToken)
	if err != nil{
		respondWithError(w, 401, fmt.Sprintf("%v", err), nil)
		return
	}
	if time.Now().After(user.ExpiresAt){
		respondWithError(w, 401, fmt.Sprintf("%v", err), nil)
		return 
	}
	if user.RevokedAt.Valid{
		respondWithError(w, 401, fmt.Sprintf("%v", err), nil)
		return
	}
	
	newToken, err := auth.MakeJWT(user.UserID, apiConfig.secret)
	if err != nil{
		respondWithError(w, 401, fmt.Sprintf("%v", err), nil)
		return
	}

	respondWithJson(w, http.StatusOK, refreshToken{
		Token:newToken,
	})
} 

func(apiConfig *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request){
	authToken, err := auth.GetBearerToken(r.Header)	
	if err != nil{
		respondWithError(w, 401, fmt.Sprint("the token is empty"), nil)
		return
	}
	err = apiConfig.db.RevokeTokenFromRefreshToken(r.Context(), authToken)
	if err != nil{
		respondWithError(w, 401, fmt.Sprintf("error revoking the token %v", err), nil)
		return
	}
	w.WriteHeader(204)
}
