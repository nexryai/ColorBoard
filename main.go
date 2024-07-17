package main

import "github.com/nexryai/ColorBoard/db"

func main() {
	client := db.NewClient()
	client.Connect()
}
