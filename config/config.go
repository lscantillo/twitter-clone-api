package config

import (
	"log"

	"github.com/joho/godotenv"
)

func GetVariables(variable string) string {
	var myEnv, err = godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return myEnv[variable]
}
