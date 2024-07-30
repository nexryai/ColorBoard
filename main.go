package main

import (
	"github.com/joho/godotenv"
	"github.com/nexryai/ColorBoard/internal/boot"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Failed loading .env file")
	}

	boot.Boot()
}
