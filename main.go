package main

import (
	"github.com/joho/godotenv"
	"github.com/nexryai/ColorBoard/db"
	"github.com/nexryai/ColorBoard/internal/boot"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Failed loading .env file")
	}

	err = db.InitializeDatabase()
	if err != nil {
		log.Print("Err: ", err)
		log.Fatalf("Failed to initialize database :(")
	}

	boot.Boot()
}
