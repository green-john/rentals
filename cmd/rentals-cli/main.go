package main

import (
	"log"
	"os"
	"rentals"
)

func main() {
	runServer()
}

func runServer() {
	addr := os.Getenv("RENTALS_ADDRESS")
	if addr == "" {
		addr = "localhost:8083"
	}

	testing := true
	if os.Getenv("RENTALS_TESTING") == "" {
		testing = false
	}

	app, err := rentals.NewApp(addr, testing)

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
