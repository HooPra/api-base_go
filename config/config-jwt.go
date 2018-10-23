package config

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

// JWTConfig for JWT
type JWTConfig struct {
	PrivateKey  string
	PublicKey   string
	ExpireDelta int
	Issuer      string
}

// JWTKeyPair is a container for this server's
// public and private JWT keys
type JWTKeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

var keyPair *JWTKeyPair

// GetJWTKeyPair returns a standard instance of a JWT key
func GetJWTKeyPair() *JWTKeyPair {
	if keyPair == nil {
		keyPair = &JWTKeyPair{
			PrivateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}

	return keyPair
}

func getPrivateKey() *rsa.PrivateKey {

	subpath := config.PrivateKeyPath
	data := readKeyFile(subpath)
	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data)

	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {

	subpath := config.PublicKeyPath
	data := readKeyFile(subpath)
	publicKeyImported, err := x509.ParsePKIXPublicKey(data)

	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		panic(err)
	}

	return rsaPub
}

func readKeyFile(subpath string) []byte {

	pwd, _ := os.Getwd()
	path := pwd + subpath
	keyFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := keyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)
	buffer := bufio.NewReader(keyFile)
	_, err = buffer.Read(pembytes)
	data, _ := pem.Decode([]byte(pembytes))

	keyFile.Close()
	return data.Bytes
}
