package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Issuer is the public name for this server
var Issuer = "restful-server.com"

var environments = map[string]string{
	"prod": "settings/prod.json",
	"dev":  "settings/dev.json",
}

// Settings for JWT
type Settings struct {
	PrivateKeyPath     string
	PublicKeyPath      string
	JWTExpirationDelta int
}

var settings Settings = Settings{}
var env = "dev"

// Init settings for the current environment
func Init() {
	env = os.Getenv("GO_ENV")
	if env == "" {
		log.Print("Warning: Using dev environment since GO_ENV was not set")
		env = "dev"
	}

	content, err := ioutil.ReadFile(environments[env])
	if err != nil {
		fmt.Println("Error while reading config file", err)
	}

	settings = Settings{}
	jsonErr := json.Unmarshal(content, &settings)
	if jsonErr != nil {
		fmt.Println("Error while parsing config file", jsonErr)
	}
}

// GetEnvironment returns the current environment
func GetEnvironment() string {
	return env
}

// Get returns the current settings
func Get() Settings {
	if &settings == nil {
		Init()
	}
	return settings
}
