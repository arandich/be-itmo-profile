package main

import (
	_ "github.com/lib/pq"
	"itmo-profile/internal/database"
	"itmo-profile/src/app"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	db.Close()

	app.Run()
}
