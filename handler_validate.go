package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: replaceBadWords(params.Body),
	})
}

func replaceBadWords(s string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	splitWord := strings.Split(s, " ")
	for i := range splitWord {
		for _, word := range badWords {
			if strings.ToLower(splitWord[i]) == word {
				splitWord[i] = "****"
			}
		}
	}
	s = strings.Join(splitWord, " ")
	return s
}
