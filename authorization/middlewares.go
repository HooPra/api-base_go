package authorization

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/hoopra/api-base_go/config"
	"github.com/hoopra/api-base_go/models"
)

// RequireTokenAuthentication ivalidates the JWT of a request
// and calls a handler if successful
func RequireTokenAuthentication(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	// printHeaders(req)
	responder := models.NewHTTPResponder(w)
	token, err := GetTokenFromRequest(req)

	if err == nil {
		valid := validateToken(token)
		if valid {
			next(w, req)
			return
		}
	}

	responder.RespondWithStatus(http.StatusUnauthorized)
}

func printHeaders(req *http.Request) {

	for name, headers := range req.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			log.Printf("%v: %v", name, h)
		}
	}
}

// GetTokenFromRequest returns a JWT from a request
// if it was signed by this server
func GetTokenFromRequest(req *http.Request) (*jwt.Token, error) {

	return request.ParseFromRequest(req, request.OAuth2Extractor, keyFunction)
}

func keyFunction(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return config.GetJWTKeyPair().PublicKey, nil
}
