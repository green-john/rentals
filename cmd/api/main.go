package main

import (
	"log"
	"tournaments"
)

func main() {
	const addr = "localhost:8083"
	app, err := rentals.NewApp(addr)

	if err != nil {
		panic(err)
	}

	// Make sure we delete all things after we are done
	log.Printf("[ERROR] %s", app.ServeHTTP())
}
