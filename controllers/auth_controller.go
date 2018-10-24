package controllers

import (
	"errors"
	"net/http"

	auth "github.com/hoopra/api-base_go/authorization"
	"github.com/hoopra/api-base_go/datastore"
	"github.com/hoopra/api-base_go/models"
)

// Register adds a user to the datastore if one
// was supplied in the request
func Register(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	responder := models.NewHTTPResponder(w)
	user := new(models.User)

	err := UnpackJSONBody(req, &user)
	if err != nil {
		responder.RespondWithError(err)
		return
	}

	if user.Username == "" || user.Password == "" {
		responder.RespondWithError(errors.New("No data supplied for user"))
		return
	}

	err = datastore.Users().Add(user)
	if err != nil {
		responder.RespondWithError(err)
		return
	}

	responder.RespondWithStatus(http.StatusAccepted)
}

// Login issues a JWT if the user in the request
// can be validated in the datastore
func Login(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	responder := models.NewHTTPResponder(w)
	user := new(models.User)
	err := UnpackJSONBody(req, &user)
	if err != nil {
		responder.RespondWithError(err)
		return
	}

	if user.Username == "" || user.Password == "" {
		responder.RespondWithError(errors.New("No credentials supplied for login"))
		return
	}

	u, err := datastore.Users().SelectByName(user.Username)
	if err != nil {
		// responder.RespondWithError(err)
		responder.RespondWithStatus(http.StatusUnauthorized)
		return
	}

	if auth.Authenticate(user) {
		token, err := auth.IssueJWT(u.UUID)
		responder.RespondWithToken(token, err)
		return
	}

	responder.RespondWithStatus(http.StatusUnauthorized)
}

// RefreshToken issues a new JWT if the request
// already contains a valid one
func RefreshToken(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	responder := models.NewHTTPResponder(w)
	user := new(models.User)
	UnpackJSONBody(req, &user)

	// decoder := json.NewDecoder(req.Body)
	// decoder.Decode(&requestUser)

	token, err := auth.IssueJWT(user.UUID)
	if err == nil {
		responder.RespondWithToken(token, err)
		return
	}

	responder.RespondWithStatus(http.StatusUnauthorized)
}

// Logout invalidates a refresh token. Dummy function for now
func Logout(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	responder := models.NewHTTPResponder(w)
	_, err := auth.GetTokenFromRequest(req)
	if err != nil {
		responder.RespondWithStatus(http.StatusInternalServerError)
	}

	responder.RespondWithStatus(http.StatusOK)
}
