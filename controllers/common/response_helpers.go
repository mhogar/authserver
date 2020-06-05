package common

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

// ParseJSONBody parses the json read by r into v
func ParseJSONBody(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(v)
	if err != nil {
		log.Println(err)
		return errors.New("invalid request body")
	}

	return nil
}

// SendResponse writes a response to w with the provided status and response.
func SendResponse(w http.ResponseWriter, status int, res interface{}) {
	//set the header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	//write the response
	encoder := json.NewEncoder(w)
	err := encoder.Encode(res)
	if err != nil {
		log.Panic(err) //panic if can't write response
	}
}

// SendSuccessResponse writes a response to w with an OK status and a basic response.
func SendSuccessResponse(w http.ResponseWriter) {
	SendResponse(w, http.StatusOK, BasicResponse{Success: true})
}

// SendErrorResponse writes a response to w with the provided status and an error response with the provided message.
func SendErrorResponse(w http.ResponseWriter, status int, messsage string) {
	SendResponse(w, status, ErrorResponse{
		Success: false,
		Error:   messsage,
	})
}

// SendInternalErrorResponse writes a response to w with InternalServerError status and an error response with an internal server error message.
func SendInternalErrorResponse(w http.ResponseWriter) {
	SendErrorResponse(w, http.StatusInternalServerError, "an internal error occurred")
}

//SendDataResponse writes a response to w with an OK status and the provided generic data.
func SendDataResponse(w http.ResponseWriter, data interface{}) {
	SendResponse(w, http.StatusOK, DataResponse{
		Success: true,
		Data:    data,
	})
}
