package authorization

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hoopra/api-base_go/config"
	"github.com/hoopra/api-base_go/datastore"
	"github.com/hoopra/api-base_go/models"
	"github.com/satori/go.uuid"
)

// IssueJWT returns a JWT signed by this server
func IssueJWT(uuid uuid.UUID) (string, error) {

	keyInstance := config.GetJWTKeyPair()
	token := jwt.New(jwt.SigningMethodRS512)

	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(config.Get().ExpireDelta)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["sub"] = uuid.String()
	claims["iss"] = config.Get().Issuer
	token.Claims = claims

	tokenString, err := token.SignedString(keyInstance.PrivateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GetUUIDFromToken returns the UUID of the user
// for which a token was issued
func GetUUIDFromToken(token *jwt.Token) uuid.UUID {

	claims := token.Claims.(jwt.MapClaims)
	idString := claims["sub"].(string)

	id, err := uuid.FromString(idString)
	if err != nil {
		panic(err)
	}

	return id
}

// Authenticate returns true if a user exists
// in the datastore
func Authenticate(user *models.User) bool {

	success := datastore.Users().Validate(user)
	return success
}

func validateToken(token *jwt.Token) bool {

	claims := token.Claims.(jwt.MapClaims)
	issuer := claims["iss"].(string)
	userID := claims["sub"].(string)
	expires := claims["exp"].(float64)
	issued := claims["iat"].(float64)

	if issuer == config.Get().Issuer && len(userID) > 0 && (expires-issued) > 0 && token.Valid {
		return true
	}

	return false
}
