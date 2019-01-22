package main

import (
	"fmt"
	"log"
	"os"
	"rentals"
)

func main() {
	runServer()
}

func runServer() {

	var addr string
	port := os.Getenv("PORT")
	if port == "" {
		addr = "localhost:8083"
	} else {
		addr = fmt.Sprintf(":%s", port)
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
