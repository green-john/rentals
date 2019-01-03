package rentals

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	// Add ID, createdAt, updatedAt, deletedAt
	gorm.Model

	// Username
	Username string `json:"username"`

	// Password hash
	PasswordHash string `json:"password_hash"`

	// Role
	Role string `json:"role"`
}

type Apartment struct {
	// Add ID, createdAt, updatedAt, deletedAt
	gorm.Model

	// Name and description
	Name, Desc string

	// Id of the realtor
	UserID string

	// Floor size area
	AreaSize float32

	// Monthly rent
	PricePerMonth float32

	// Number of rooms
	RoomCount int

	// Geolocation
	Lat, Long float32
}
