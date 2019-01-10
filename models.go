package rentals

// TODO add createdAt

type User struct {
	// Primary key
	ID uint `gorm:"primary_key",json:"id"`

	// Username
	Username string `json:"username"`

	// Password hash. Not included in json responses
	PasswordHash string

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

type Apartment struct {
	// Primary key
	ID uint `gorm:"primary_key",json:"id"`

	// Name of this property
	Name string `json:"name"`

	// Description
	Desc string `json:"description"`

	// Realtor associated with this apartment
	Realtor   User `gorm:"foreignkey:RealtorId"`
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

// Validates data for a new apartment.
func (s *Apartment) Validate() error {
	return nil
}

var DbModels = []interface{}{
	&User{},
	&UserSession{},
	&Apartment{},
}
