package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	env := os.Getenv("APP_ENV")
	switch env {
	case "remote":
		if err := godotenv.Load(".env.remote"); err != nil {
			log.Println("no .env.remote found, using system env vars")
		} else {
			log.Println("loaded .env.remote")
		}
	default:
		if err := godotenv.Load(".env"); err != nil {
			log.Println("no .env found, using system env vars")
		} else {
			log.Println("loaded .env")
		}
	}
}
