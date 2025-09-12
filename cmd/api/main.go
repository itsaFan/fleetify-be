package main

import (
	"log"
	"net/http"

	"github.com/itsaFan/fleetify-be/internal/config"
	apihttp "github.com/itsaFan/fleetify-be/internal/http"
)

func main() {
	config.LoadEnv()
	db, err := config.DBConnection()
	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}

	router := apihttp.NewRouter(db)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("listening on port 8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
