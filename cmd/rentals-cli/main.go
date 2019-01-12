package main

import (
	"fmt"
	"log"
	"tournaments"
)

func main() {
	createAdmin("admin", "admin")
	runServer()
}

func createAdmin(username, password string) {
	db, err := rentals.ConnectToDB()

	if err != nil {
		panic(err)
	}

	ur := &rentals.UserResource{Db: db}

	jsonString := fmt.Sprintf(`{"username": "%s", "password": "%s", "role": "admin"}`,
		username, password)
	_, err = ur.Create([]byte(jsonString))

	if err != nil {
		panic(err)
	}
}

func runServer() {
	const addr = "localhost:8083"
	app, err := rentals.NewApp(addr)

	if err != nil {
		panic(err)
	}

	err = app.Setup()
	if err != nil {
		panic(err)
	}

	// Make sure we delete all things after we are done
	log.Printf("[ERROR] %s", app.ServeHTTP())
}
