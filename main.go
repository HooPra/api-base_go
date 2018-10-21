package main

import (
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/hoopra/api-base_go/routing"
	"github.com/hoopra/api-base_go/settings"
	"github.com/rs/cors"
)

func main() {

	// Set up environment
	os.Setenv("GO_ENV", "dev")
	settings.Init()

	// Create router
	router := routing.GetRouting()
	n := negroni.Classic()

	corsOpts := cors.Options{}
	corsOpts.AllowCredentials = true
	corsOpts.AllowedOrigins = []string{"*"}
	corsOpts.AllowedHeaders = []string{"Origin", "Content-Type", "Authorization"}

	handler := cors.New(corsOpts).Handler(router)
	n.UseHandler(handler)

	// Run server
	log.Print("Listening on " + routing.Port)
	http.ListenAndServe(routing.Port, n)
}

func corsHandler() {

}
