package main

import (
	"fmt"
	"log"
	"os"
	"rentals"
	"rentals/services"
	"rentals/transport"
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

	db, err := rentals.ConnectToDB(testing)
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(rentals.DbModels...)

	authN := services.NewDbAuthnService(db)
	authZ := services.NewAuthzService()
	apartmentsSrv := services.NewDbApartmentService(db)
	userService := services.NewDbUserService(db)

	srv, err := transport.NewServer(db, authN, authZ, apartmentsSrv, userService)
	if err != nil {
		log.Fatal(err)
	}

	// Make sure we delete all things after we are done
	log.Printf("[ERROR] %s", srv.ServeHTTP(addr))
}
