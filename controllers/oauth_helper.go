package controllers

import (
	requesterror "authserver/common/request_error"
	commonhelpers "authserver/helpers/common"
	"authserver/models"
	"log"

	"github.com/google/uuid"
)

func parseClient(clientCRUD models.ClientCRUD, clientID uuid.UUID) (*models.Client, requesterror.OAuthRequestError) {
	//get the client
	client, err := clientCRUD.GetClientByID(clientID)
	if err != nil {
		log.Println(commonhelpers.ChainError("error getting client by id", err))
		return nil, requesterror.OAuthInternalError()
	}

	//check client was found
	if client == nil {
		return nil, requesterror.OAuthClientError("invalid_client", "client with id not found")
	}

	return client, requesterror.OAuthNoError()
}

func parseScope(scopeCRUD models.ScopeCRUD, name string) (*models.Scope, requesterror.OAuthRequestError) {
	//get the scope
	scope, err := scopeCRUD.GetScopeByName(name)
	if err != nil {
		log.Println(commonhelpers.ChainError("error getting scope by name", err))
		return nil, requesterror.OAuthInternalError()
	}

	if scope == nil {
		return nil, requesterror.OAuthClientError("invalid_scope", "scope with name not found")
	}

	return scope, requesterror.OAuthNoError()
}
