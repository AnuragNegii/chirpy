package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AnuragNegii/chirpy/internal/auth"
	"github.com/AnuragNegii/chirpy/internal/database"
	"github.com/google/uuid"
)

type Posts struct{
	ID uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Body string `json:"body"`
	User_id uuid.UUID `json:"user_id"`
}

func (apiConfig *apiConfig) hadnleChirps(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body) 
	var params parameters
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength{
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	strList := strings.Split(params.Body, " ")
	wordsList := []string{"kerfuffle", "sharbert", "fornax"}
	for i, word := range strList{
		for _, badWord := range wordsList{
			if strings.ToLower(word) == badWord{
				strList[i] = "****"
			}
		}
	}
	jwtString, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w ,http.StatusBadRequest, fmt.Sprintf("%v",err), nil)
		return
	}
	userID, err := auth.ValidateJWT(jwtString, apiConfig.secret)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("%v", err), nil)
		return 
	}
	
	params.Body = strings.Join(strList, " ")
	post, err := apiConfig.db.CreateChirps(r.Context(), database.CreateChirpsParams{
		Body: params.Body,
		UserID: userID,
	})
	if err != nil{
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("not able to create chirps: %v", err), nil)
		return
	}
	respondWithJson(w, http.StatusCreated, Posts{
			ID: post.ID,
			Created_at: post.CreatedAt,
			Updated_at: post.UpdatedAt,
			Body: post.Body,
			User_id: userID,
		})
}

func (apiConfig *apiConfig) GetChirps(w http.ResponseWriter, r *http.Request){
	var chirps []Posts
	arrChirps, err := apiConfig.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("cant retrieve all chirps: %v", err), nil)
		return
	}
	for _, r := range arrChirps{
		chirps = append(chirps, Posts{
			ID: r.ID,
			Created_at: r.CreatedAt,
			Updated_at: r.UpdatedAt,
			Body: r.Body,
			User_id: r.UserID,
		})
	}
	respondWithJson(w, http.StatusOK, chirps)
}

func (apiConfig *apiConfig) GetChirpsById(w http.ResponseWriter, r *http.Request){
	chirpId := r.PathValue("chirpID")
	uuidID, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("cant parse chirpID: %v", err), nil)
		return
	}
	chirp, err := apiConfig.db.GetChirpByID(r.Context(), uuidID)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("cant get that chirp: %v", err), nil)
		return
	}
	respondWithJson(w, http.StatusOK, Posts{
		ID: chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body: chirp.Body,
		User_id: chirp.UserID,
	})
}
