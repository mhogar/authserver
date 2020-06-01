package controllers

import (
	"log"
	"net/http"
)

// BasicResponse represents a response with a simple true/false success field
type BasicResponse struct {
	Success bool `json:"success"`
}

// ErrorResponse represents a response with a true/false success field and an error message
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// DataResponse represents a response with a true/false success field and generic data
type DataResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// PanicHandler is the function to be called if a panic is encountered
func PanicHandler(w http.ResponseWriter, req *http.Request, info interface{}) {
	log.Println(info)

	sendInternalErrorResponse(w)
}
