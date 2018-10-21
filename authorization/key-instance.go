package authorization

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hoopra/api-base_go/datastore"
	"github.com/hoopra/api-base_go/models"
	"github.com/hoopra/api-base_go/settings"
	"github.com/satori/go.uuid"
)

// JWTKeyInstance is a container for this server's
// public and private JWT keys
type jwtKeyInstance struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

const (
	expireOffset = 3600
)

var keyInstance *jwtKeyInstance

// GetJWTKeyInstance returns a standard instance of a JWT key
func getJWTKeyInstance() *jwtKeyInstance {
	if keyInstance == nil {
		keyInstance = &jwtKeyInstance{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}

	return keyInstance
}

// GenerateAccessToken returns a JWT signed by this server
func GenerateAccessToken(uuid uuid.UUID) (string, error) {

	keyInstance := getJWTKeyInstance()
	token := jwt.New(jwt.SigningMethodRS512)

	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["sub"] = uuid.String()
	claims["iss"] = settings.Issuer
	token.Claims = claims

	tokenString, err := token.SignedString(keyInstance.privateKey)
	if err != nil {
		// panic(err)
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

	success := datastore.Store().Users().Validate(user)
	return success
}

func getTokenRemainingValidity(timestamp interface{}) int {

	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}

	return expireOffset
}

func validateToken(token *jwt.Token) bool {

	claims := token.Claims.(jwt.MapClaims)
	issuer := claims["iss"].(string)
	userID := claims["sub"].(string)
	expires := claims["exp"].(float64)
	issued := claims["iat"].(float64)

	if issuer == settings.Issuer && len(userID) > 0 && (expires-issued) > 0 && token.Valid {
		return true
	}

	return false
}

func getPrivateKey() *rsa.PrivateKey {

	pwd, _ := os.Getwd()
	path := pwd + settings.Get().PrivateKeyPath

	privateKeyFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {

	pwd, _ := os.Getwd()
	path := pwd + settings.Get().PublicKeyPath
	publicKeyFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	publicKeyFile.Close()

	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		panic(err)
	}

	return rsaPub
}
