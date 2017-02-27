package authorization

import (
	"fmt"
	"net/http"
	"strings"

	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/hoopra/GoAuthServer/models"
)

// RequireTokenAuthentication is a middleware component that
// validates the JWT of a request and calls a secured handler
// function if successful
func RequireTokenAuthentication(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	// Loop through headers
	for name, headers := range req.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			log.Printf("%v: %v", name, h)
		}
	}

	responder := models.NewHTTPResponder(w)
	keyInstance := GetJWTKeyInstance()
	token, err := GetTokenFromRequest(req)

	if err == nil {
		valid := keyInstance.validateToken(token)
		if valid {
			next(w, req)
			return
		}
		return
	}

	responder.RespondWithStatus(http.StatusUnauthorized)
}

// GetTokenFromRequest returns a JWT from a request
// if it was signed by this server
func GetTokenFromRequest(req *http.Request) (*jwt.Token, error) {

	return request.ParseFromRequest(req, request.OAuth2Extractor, keyFunction)
}

func keyFunction(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	} else {
		return keyInstance.PublicKey, nil
	}
}
