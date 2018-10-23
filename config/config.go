package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var environments = map[string]string{
	"prod": "config/prod.json",
	"dev":  "config/dev.json",
}

// Config for instance
type Config struct {
	PrivateKeyPath string
	PublicKeyPath  string
	ExpireDelta    int
	Issuer         string
	Env            string
}

var config = Config{}
var env string

// Init Config for the current environment
func Init() {
	env = os.Getenv("GO_ENV")
	if env == "" {
		log.Print("using dev as GO_ENV by default")
		env = "dev"
	}

	// fmt.Printf("Threads: %v\n", runtime.GOMAXPROCS(-1))

	content, err := ioutil.ReadFile(environments[env])
	if err != nil {
		fmt.Println("error while reading config file", err)
	}

	jsonErr := json.Unmarshal(content, &config)
	if jsonErr != nil {
		fmt.Println("error while parsing config file", jsonErr)
	}

	priOvr := os.Getenv("AUTH_PRIVATE_KEY_PATH")
	pubOvr := os.Getenv("AUTH_PUBLIC_KEY_PATH")
	if priOvr != "" && pubOvr != "" {
		config.PrivateKeyPath = priOvr
		config.PublicKeyPath = pubOvr
	}

	config.Env = env
}

// GetConfig returns the current config settings
func Get() Config {
	if &config == nil {
		Init()
	}
	return config
}
