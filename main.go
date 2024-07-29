package main

import (
	"context"
	"github.com/nexryai/ColorBoard/db"
)

func main() {
	client := db.NewClient()
	client.Connect()

	user, err := client.User.FindFirst(db.User.ID.Equals("1")).Exec(context.TODO())
	if user == nil {
		println("User not found")
	}

	println(user.ID, err)
}
