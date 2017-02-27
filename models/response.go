package models

import (
	"encoding/json"
	"net/http"
)

// HTTPResponder wraps an HTTP response writer
type HTTPResponder struct {
	writer http.ResponseWriter
}

// NewHTTPResponder returns an HTTPResponder
func NewHTTPResponder(w http.ResponseWriter) *HTTPResponder {
	return &HTTPResponder{writer: w}
}

type accessTokenResponse struct {
	Token string `json:"token" form:"token"`
}

type textResponse struct {
	Text string `json:"text"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// RespondWithError writes a JSON text response to the client
func (responder *HTTPResponder) RespondWithJSONText(text string) {

	responder.setType("application/json")
	response := textResponse{Text: text}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		responder.RespondWithStatus(http.StatusInternalServerError)
	} else {
		responder.respondWithBody(http.StatusOK, jsonResponse)
	}
}

// RespondWithError writes an access token to the client
func (responder *HTTPResponder) RespondWithToken(token string, err error) {

	responder.setType("application/json")
	if err != nil {
		responder.RespondWithStatus(http.StatusInternalServerError)
	} else {
		response, _ := json.Marshal(accessTokenResponse{token})
		responder.respondWithBody(http.StatusOK, response)
	}
}

// RespondWithError writes an error to the client
func (responder HTTPResponder) RespondWithError(err error) {

	responder.setType("application/json")
	response := errorResponse{Error: err.Error()}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		responder.RespondWithStatus(http.StatusInternalServerError)
	} else {
		responder.respondWithBody(http.StatusBadRequest, jsonResponse)
	}
}

func (responder *HTTPResponder) setType(typeString string) {
	responder.writer.Header().Set("Content-Type", typeString)
}

// RespondWithStatus writes an http status code
// with an empty body to the client
func (responder *HTTPResponder) RespondWithStatus(status int) {
	responder.respondWithBody(status, []byte(""))
}

func (responder *HTTPResponder) respondWithBody(status int, body []byte) {
	// responder.writer.Header().Add("Access-Control-Allow-Origin", "http://localhost:4200")
	// responder.writer.Header().Add("Access-Control-Allow-Methods", "GET, POST")
	responder.writer.WriteHeader(status)
	responder.writer.Write(body)
}
