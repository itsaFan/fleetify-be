package main

import (
	"log"
	"net/http"

	// "github.com/itsaFan/fleetify-be/internal/config"
	// apihttp "github.com/itsaFan/fleetify-be/internal/http"
	// "github.com/itsaFan/fleetify-be/internal/model"
)

func main() {

	srv := &http.Server{
		Addr: ":8080",
		// Handler: router,
	}

	log.Println("listening on port 8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
