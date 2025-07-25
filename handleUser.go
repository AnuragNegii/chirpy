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
	Token string `json:"token"`
	Refresh_token string `json:"refresh_token"`
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

func (apiConfig *apiConfig) changeUserEmailAndPass(w http.ResponseWriter, r *http.Request){
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, 401, fmt.Sprintf("response token missing %v", err), nil)
		return
	}
	
	userId, err := auth.ValidateJWT(tokenString, apiConfig.secret)
	if err != nil{
		respondWithError(w, 401, fmt.Sprintf("not a validJWT %v", err), nil)
		return
	}

	type EmailBody struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	defer r.Body.Close()
	var emailBody EmailBody
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&emailBody)
	if err != nil{
		respondWithError(w, 401, fmt.Sprintf("%v", err), nil)
		return
	}
	user, err := apiConfig.db.GetUserFromUserID(r.Context(), userId)
	if err != nil{
		respondWithError(w, 401, fmt.Sprintf("%v", err), nil)
		return
	}

	hashPass, err := auth.HashPassword(emailBody.Password)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err), nil)
		return
	}
	err = apiConfig.db.ChangeUserNameAndPassword(r.Context(), database.ChangeUserNameAndPasswordParams{
		Email: emailBody.Email,
		HashedPassword: hashPass,
		ID: user.ID,
	})	
	if err != nil{
		respondWithError(w, 401, fmt.Sprintf("%v", err), nil)
		return
	}
	respondWithJson(w, http.StatusOK, User{
		ID: user.ID,	
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,	
		Email: emailBody.Email,
		Token: tokenString,
	})
}
