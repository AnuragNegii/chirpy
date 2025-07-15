package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/AnuragNegii/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirps struct{
	Id uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Body string `json:"body"`
	User_id uuid.UUID `json:"user_id"`
}

func (apiConfig *apiConfig) handleChirp(w http.ResponseWriter, r *http.Request){
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

	params.Body = strings.Join(strList, " ")
	user, err := apiConfig.db.CreateChirps(context.Background(), database.CreateChirpsParams{
		UserID: params.User_id,
		Body: params.Body,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error creating chirp in db", nil)
		return
	}
	respondWithJson(w, http.StatusOK, Chirps{
		Id: user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Body: user.Body,
		User_id: user.ID,
	})
}
