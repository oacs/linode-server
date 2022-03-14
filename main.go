package main

import (
	"os"

	"example.com/m/v2/modules/server"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	args := os.Args
	envDir := ".env"
	if len(args) > 1 {
		if args[1] == "--prod" {
			envDir = "/home/www/.env"
		}
	}

	err := godotenv.Load(envDir)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	server.Main()
}
