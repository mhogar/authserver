package router

import (
	"log"
	"net/http"
)

// panicHandler is the function to be called if a panic is encountered
func panicHandler(w http.ResponseWriter, req *http.Request, info interface{}) {
	log.Println(info)
	sendInternalErrorResponse(w)
}
