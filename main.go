package main

import (
	"example.com/m/v2/modules/server"
	log "github.com/sirupsen/logrus"
	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	server.Main()
}
