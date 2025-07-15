package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

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
		User_id uuid.UUID `json:"user_id"`
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

	post, err := apiConfig.db.CreateChirps(context.Background(), database.CreateChirpsParams{
		Body: params.Body,
		UserID: params.User_id,
	})
	if err != nil{
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("not able to create chirps: %v", err), nil)
		return
	}

	params.Body = strings.Join(strList, " ")
	respondWithJson(w, http.StatusCreated, Posts{
			ID: post.ID,
			Created_at: post.CreatedAt,
			Updated_at: post.UpdatedAt,
			Body: post.Body,
			User_id: post.UserID,
		})
}
