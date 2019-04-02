package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

// Creates a new connection to the Db and migrates
// all the objects so that tables in the Db are created
func ConnectToDB(testing bool) (*gorm.DB, error) {
	const dialect = "postgres"

	extras := ""
	if testing {
		extras += "sslmode=disable"
	}

	dbHost := os.Getenv("RENTALS_DB_HOST")
	dbName := os.Getenv("RENTALS_DB_NAME")
	dbUser := os.Getenv("RENTALS_DB_USER")
	dbPass := os.Getenv("RENTALS_DB_PASSWORD")

	cnxString := fmt.Sprintf("host=%s user=%s dbname=%s", dbHost, dbUser, dbName)

	if dbPass != "" {
		cnxString += fmt.Sprintf(" password=%s", dbPass)
	}

	if extras != "" {
		cnxString += " " + extras
	}

	db, err := gorm.Open(dialect, cnxString)
	if err != nil {
		return nil, fmt.Errorf("[ConnectToDB] error calling gorm.Open(): %v", err)
	}

	return db, nil
}
