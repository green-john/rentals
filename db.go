package rentals

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
)

// Creates a new connection to the Db and migrates
// all the objects so that tables in the db are created
func ConnectToDB() (*gorm.DB, error) {
	const dialect = "postgres"
	const extras = "sslmode=disable"

	dbHost := os.Getenv("RENTALS_DB_HOST")
	dbName := os.Getenv("RENTALS_DB_NAME")
	dbUser := os.Getenv("RENTALS_DB_USER")

	cnxString := fmt.Sprintf("host=%s user=%s dbname=%s %s", dbHost, dbUser, dbName, extras)
	log.Printf("Connecting to [%s]", cnxString)

	db, err := gorm.Open(dialect, cnxString)
	if err != nil {
		return nil, fmt.Errorf("[ConnectToDB] error calling gorm.Open(): %v", err)
	}

	log.Print("Connected")


	return db, nil
}
