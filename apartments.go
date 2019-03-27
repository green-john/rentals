package rentals

import (
	"errors"
	"time"
)

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

type ApartmentService interface {
	Create(ApartmentCreateInput) (*ApartmentCreateOutput, error)
	Read(ApartmentReadInput) (*ApartmentReadOutput, error)
	Find(ApartmentFindInput) (*ApartmentFindOutput, error)
	Update(ApartmentUpdateInput) (*ApartmentUpdateOutput, error)
	Delete(ApartmentDeleteInput) (*ApartmentDeleteOutput, error)
}

type ApartmentCreateInput struct {
	Apartment
}

type ApartmentCreateOutput struct {
	Apartment
}

type ApartmentReadInput struct {
	// ID to lookup the apartment
	Id string
}

type ApartmentReadOutput struct {
	Apartment
}

type ApartmentFindInput struct {
	Query string
}

type ApartmentFindOutput struct {
	Apartments []Apartment
}

func (o *ApartmentFindOutput) Public() interface{} {
	return o.Apartments
}

type ApartmentUpdateInput struct {
	Id   string
	Data map[string]interface{}
}

type ApartmentUpdateOutput struct {
	Apartment
}

type ApartmentDeleteInput struct {
	Id string
}

type ApartmentDeleteOutput struct {
	Message string
}
