package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KalimaraPeleteiro/RSS-Aggregator/internal/database"
	"github.com/google/uuid"
)

// Arquivo responsável por lidar com as requisições

// Rotas padrão para verificação
func handlerDefaultError(response http.ResponseWriter, request *http.Request) {
	errorJSON(response, 400, "Esta é uma mensagem padrão de erro.")
}

func handlerServerReadiness(response http.ResponseWriter, request *http.Request) {
	JSONResponse(response, 200, struct{}{})
}

// Adicionando Feeds aos Usuários
func (apiConfiguration *apiConfig) handlerCreateFeed(response http.ResponseWriter, request *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"nome"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorJSON(response, 400, fmt.Sprintf("Erro ao decodificar JSON: %s", err))
		return
	}

	feed, err := apiConfiguration.database.CreateFeed(request.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})

	if err != nil {
		errorJSON(response, 400, fmt.Sprintf("Erro ao adicionar feed: %s", err))
		return
	}

	JSONResponse(response, 201, SQLCFeedToFeed(feed))

}

// Retornando todos os Feeds
func (apiConfiguration apiConfig) handlerGetAllFeeds(response http.ResponseWriter, request *http.Request, user database.User) {
	feeds, err := apiConfiguration.database.GetFeeds(request.Context(), user.ID)
	if err != nil {
		errorJSON(response, 400, fmt.Sprintf("Não consegui buscar os feeds: %v", err))
		return
	}

	JSONResponse(response, 200, feeds)
}
