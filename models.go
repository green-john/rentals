package rentals

import (
	"errors"
	"time"
)

type uid uint

type User struct {
	// Primary key
	ID uid `gorm:"primary_key" json:"id"`

	// Username
	Username string `gorm:"unique" json:"username"`

	// Password hash. Not included in json responses
	PasswordHash string `json:"-"`

	// Role
	Role string `json:"role"`
}

type UserSession struct {
	// Primary key
	ID uint `gorm:"primary_key"`

	// Generated token
	Token string

	// User associated to this session
	UserID uint
	User   User
}

// TODO add createdAt
type Apartment struct {
	// Primary key
	ID uid `gorm:"primary_key" json:"id"`

	// Date added
	DateAdded time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"dateAdded"`

	// Name of this property
	Name string `json:"name"`

	// Description
	Desc string `json:"description"`

	// Realtor associated with this apartment
	Realtor   User `json:"-" gorm:"foreignkey:RealtorId"`
	RealtorId uint `json:"realtorId"`

	// Floor size area
	// See:
	// https://stackoverflow.com/questions/445191/should-we-put-units-of-measurements-in-attribute-names
	FloorAreaMeters float32 `json:"floorAreaMeters"`

	// Monthly rent
	PricePerMonthUsd float32 `json:"pricePerMonthUSD"`

	// Number of rooms
	RoomCount int `json:"roomCount"`

	// Geolocation
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`

	// Availability of the apartment
	Available bool `json:"available"`
}

func (uid) UnmarshalJSON([]byte) error {
	return nil
}

// Validates data for a new apartment.
func (s *Apartment) Validate() error {
	allErrors := ""

	if s.Name == "" {
		allErrors += "Name can't be empty\n"
	}

	if s.FloorAreaMeters <= 0 {
		allErrors += "Floor Area must be greater than 0\n"
	}

	if s.PricePerMonthUsd <= 0 {
		allErrors += "Price per month must be greater than 0\n"
	}

	if s.RoomCount <= 0 {
		allErrors += "Room count must be greater than 0\n"
	}

	if s.Latitude < -90 || s.Latitude > 90 {
		allErrors += "Latitude must be in the range [-90.0, 90.0]\n"
	}

	if s.Longitude < -180 || s.Longitude > 180 {
		allErrors += "Longitude must be in the range [-180.0, 180.0]\n"
	}

	if allErrors == "" {
		return nil
	}

	return errors.New("\n" + allErrors)
}

var DbModels = []interface{}{
	&User{},
	&UserSession{},
	&Apartment{},
}
