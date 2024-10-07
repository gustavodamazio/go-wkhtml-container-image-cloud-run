package main

import (
	"context"
	"log"
	"net/http"
	"os"

	gorilaHandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gustavodamazio/go-test/handlers"
	"github.com/gustavodamazio/go-test/middlewares"
	"github.com/gustavodamazio/go-test/services/storage"
)

func main() {
	// #region Initialize Google Cloud Storage service
	bucketName := os.Getenv("GCS_BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("Environment variable GCS_BUCKET_NAME is not set")
	}
	ctx := context.Background()
	storageService, err := storage.NewStorageService(ctx, bucketName)
	if err != nil {
		log.Fatalf("Failed to initialize storage service: %v", err)
	}
	defer storageService.Close()
	// #endregion

	r := mux.NewRouter()
	r.Use(gorilaHandler.RecoveryHandler())
	r.Use(middlewares.RequireJSONMiddleware)
	r.Use(middlewares.RequirePOSTMiddleware)

	r.HandleFunc("/", handlers.HandleHtmlToPDF(storageService))

	loggedRouter := gorilaHandler.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))
}
