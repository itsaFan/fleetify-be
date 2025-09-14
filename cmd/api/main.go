package main

import (
	"log"
	"net/http"
	"os"
	_ "time/tzdata"

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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	log.Println("listening on port", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
