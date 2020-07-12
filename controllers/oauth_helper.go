package controllers

import (
	commonhelpers "authserver/helpers/common"
	"authserver/models"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func sendOAuthErrorResponse(w http.ResponseWriter, status int, err string, description string) {
	sendResponse(w, status, OAuthErrorResponse{
		Error:            err,
		ErrorDescription: description,
	})
}

func parseClient(clientCRUD models.ClientCRUD, w http.ResponseWriter, clientIDStr string) *models.Client {
	//parse the client id
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		log.Println(commonhelpers.ChainError("error parsing client id", err))
		sendOAuthErrorResponse(w, http.StatusUnauthorized, "invalid_client", "client_id was in invalid format")
		return nil
	}

	//get the client
	client, err := clientCRUD.GetClientByID(clientID)
	if err != nil {
		log.Println(commonhelpers.ChainError("error getting client by id", err))
		sendInternalErrorResponse(w)
		return nil
	}

	if client == nil {
		sendOAuthErrorResponse(w, http.StatusUnauthorized, "invalid_client", "")
		return nil
	}

	return client
}

func parseScope(scopeCRUD models.ScopeCRUD, w http.ResponseWriter, name string) *models.Scope {
	//get the scope
	scope, err := scopeCRUD.GetScopeByName(name)
	if err != nil {
		log.Println(commonhelpers.ChainError("error getting scope by name", err))
		sendInternalErrorResponse(w)
		return nil
	}

	if scope == nil {
		sendOAuthErrorResponse(w, http.StatusBadRequest, "invalid_scope", "")
		return nil
	}

	return scope
}

func parseAuthHeader(accessTokenCRUD models.AccessTokenCRUD, w http.ResponseWriter, req *http.Request) *models.AccessToken {
	//extract the token string from the authorization header
	splitTokens := strings.Split(req.Header.Get("Authorization"), "Bearer ")
	if len(splitTokens) != 2 {
		sendErrorResponse(w, http.StatusUnauthorized, "no bearer token provided")
		return nil
	}

	//parse the token
	tokenID, err := uuid.Parse(splitTokens[1])
	if err != nil {
		log.Println(commonhelpers.ChainError("error parsing access token id", err))
		sendErrorResponse(w, http.StatusUnauthorized, "bearer token was in invalid format")
		return nil
	}

	//fetch the token
	token, err := accessTokenCRUD.GetAccessTokenByID(tokenID)
	if err != nil {
		log.Println(commonhelpers.ChainError("error getting access token by id", err))
		sendInternalErrorResponse(w)
		return nil
	}

	if token == nil {
		sendErrorResponse(w, http.StatusUnauthorized, "invalid bearer token")
	}

	return token
}
