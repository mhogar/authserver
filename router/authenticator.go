package router

import (
	requesterror "authserver/common/request_error"
	"authserver/models"
	"net/http"
)

type Authenticator interface {
	Authenticate(req *http.Request) (*models.AccessToken, requesterror.RequestError)
}
