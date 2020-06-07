package controllers

import (
	"authserver/database"
	"authserver/helpers"
	"authserver/models"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type passwordGrantBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ClientID string `json:"client_id"`
	Scope    string `json:"scope"`
}

func sendOAuthErrorResponse(w http.ResponseWriter, status int, err string, description string) {
	sendResponse(w, status, OAuthErrorResponse{
		Error:            err,
		ErrorDescription: description,
	})
}

func parseClient(clientCRUD database.ClientCRUD, w http.ResponseWriter, clientIDStr string) *models.Client {
	//parse the client id
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		log.Println(helpers.ChainError("error parsing client id", err))
		sendOAuthErrorResponse(w, http.StatusUnauthorized, "invalid_client", "client_id was in invalid format")
		return nil
	}

	//get the client
	client, err := clientCRUD.GetClientByID(clientID)
	if err != nil {
		log.Println(helpers.ChainError("error getting client by id", err))
		sendInternalErrorResponse(w)
		return nil
	}

	if client == nil {
		sendOAuthErrorResponse(w, http.StatusUnauthorized, "invalid_client", "")
		return nil
	}

	return client
}

func parseScope(scopeCRUD database.ScopeCRUD, w http.ResponseWriter, name string) *models.Scope {
	//get the scope
	scope, err := scopeCRUD.GetScopeByName(name)
	if err != nil {
		log.Println(helpers.ChainError("error getting scope by name", err))
		sendInternalErrorResponse(w)
		return nil
	}

	if scope == nil {
		sendOAuthErrorResponse(w, http.StatusUnauthorized, "invalid_scope", "")
		return nil
	}

	return scope
}
