package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"rentals"
	"rentals/services"
	"rentals/transport"
	"strconv"
)

func main() {
	testing := flag.Bool("local", false, "runs the server with a local db")
	port := flag.Int("port", 8083, "port to bind to")

	flag.Parse()

	runServer(*testing, *port)
}

func runServer(testing bool, port int) {
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
		_, _ = fmt.Fprintln(os.Stderr, "error creating server")
		os.Exit(1)
	}

	portStr := os.Getenv("PORT")
	if portStr != "" {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error parsing port: %s\n", portStr)
			os.Exit(1)
		}
	}
	addr := fmt.Sprintf(":%d", port)
	_, _ = fmt.Fprintf(os.Stderr, "Running in %s\n", addr)
	_, _ = fmt.Fprintf(os.Stderr, "[ERROR] %s", srv.ServeHTTP(addr))
}
