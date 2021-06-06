package router

import (
	"authserver/common"
	requesterror "authserver/common/request_error"
	"authserver/database"
	"authserver/models"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type handlerFunc func(req *http.Request, params httprouter.Params, token *models.AccessToken, tx database.Transaction) (int, interface{})

func (h RouterFactory) createHandler(handler handlerFunc, authenticateUser bool) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		var token *models.AccessToken
		var rerr requesterror.RequestError

		//authenticate the user if required
		if authenticateUser {
			token, rerr = h.Authenticator.Authenticate(req)
			if rerr.Type == requesterror.ErrorTypeClient {
				sendErrorResponse(w, http.StatusUnauthorized, rerr.Error())
				return
			} else if rerr.Type == requesterror.ErrorTypeInternal {
				sendInternalErrorResponse(w)
				return
			}
		}

		//start a new transaction
		tx, err := h.TransactionFactory.CreateTransaction()
		if err != nil {
			log.Println(common.ChainError("error creating transaction", err))
			sendInternalErrorResponse(w)
			return
		}

		//execute the handler, commit the transaction on success, rollback on error
		status, body := handler(req, params, token, tx)
		if status == http.StatusOK {
			//commit the transaction
			err = tx.CommitTransaction()
			if err != nil {
				log.Println(common.ChainError("error commiting transaction", err))
				sendInternalErrorResponse(w)
				return
			}
		} else {
			tx.RollbackTransaction()
		}

		//send the response
		sendResponse(w, status, body)
	}
}
