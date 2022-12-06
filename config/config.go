package config

import (
	"log"

	"github.com/joho/godotenv"
)

func InitConfig() {
	if envErr := godotenv.Load(); envErr != nil {
		log.Fatal("Error loading .env file")
	}
}
