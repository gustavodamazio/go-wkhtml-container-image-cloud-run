package main

import (
	"log"
	"net/http"
	"os"

	gorilaHandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gustavodamazio/go-test/handlers"
	"github.com/gustavodamazio/go-test/middlewares"
)

func main() {
	r := mux.NewRouter()
	r.Use(gorilaHandler.RecoveryHandler())
	r.Use(middlewares.RequireJSONMiddleware)
	r.Use(middlewares.RequirePOSTMiddleware)

	r.HandleFunc("/", handlers.HandleHtmlToPDF)

	loggedRouter := gorilaHandler.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))
}
