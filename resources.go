package rentals

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

type UserResource struct {
	Db *gorm.DB
}

func (t *UserResource) Name() string {
	return "users"
}

func (t *UserResource) Create(jsonData []byte) ([]byte, error) {
	var newUserSchema struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	err := json.Unmarshal(jsonData, &newUserSchema)
	if err != nil {
		return nil, fmt.Errorf("[UserResource.Create] error calling json.Unmarshall(): %v", err)
	}

	if !contains([]string{"admin", "realtor", "client"}, newUserSchema.Role) {
		return nil, errors.New(
			fmt.Sprintf("[UserResource.Create] error creating user. Unknown role %s", newUserSchema.Role))
	}

	pwdHash, err := EncryptPassword(newUserSchema.Password)
	if err != nil {
		return nil, fmt.Errorf("[UserResource.Create] error encrypting password %v", err)
	}

	user := User{
		Username:     newUserSchema.Username,
		PasswordHash: pwdHash,
		Role:         newUserSchema.Role,
	}

	t.Db.Create(&user)

	rawJson, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	return rawJson, nil
}

func (t *UserResource) Read(id string) ([]byte, error) {
	panic("implement me")
}

func (t *UserResource) Update(id string, jsonData []byte) ([]byte, error) {
	panic("implement me")
}

func (t *UserResource) Delete(id string) error {
	panic("implement me")
}

type ApartmentResource struct {
	Db *gorm.DB
}

func (r *ApartmentResource) Name() string {
	return "apartments"
}

func (r *ApartmentResource) Create(jsonData []byte) ([]byte, error) {
	var newApartment NewApartmentSchema

	err := json.Unmarshal(jsonData, &newApartment)
	if err != nil {
		return nil, fmt.Errorf("[ApartmentResource.Create] error calling json.Unmarshall(): %v", err)
	}

	realtor := r.getUser(newApartment.RealtorId)
	if realtor == nil {
		return nil, fmt.Errorf("realtor %d not found", newApartment.RealtorId)
	}

	err = newApartment.Validate()
	if err != nil {
		return nil, err
	}

	apartment := Apartment{
		Name:             newApartment.Name,
		Desc:             newApartment.Desc,
		RealtorId:        newApartment.RealtorId,
		Realtor:          *realtor,
		FloorAreaMeters:  newApartment.FloorAreaMeters,
		PricePerMonthUsd: newApartment.PricePerMonthUsd,
		RoomCount:        newApartment.RoomCount,
		Latitude:         newApartment.Latitude,
		Longitude:        newApartment.Longitude,
		Available:        newApartment.Available,
	}

	r.Db.Create(&apartment)

	rawJson, err := json.Marshal(apartment)
	if err != nil {
		return nil, err
	}

	return rawJson, nil
}

func (r *ApartmentResource) getUser(userId uint) *User {
	var user User
	r.Db.Where("id = ?", userId).First(&user)

	if user.ID == 0 {
		return nil
	}

	return &user
}

func (r *ApartmentResource) Read(id string) ([]byte, error) {
	panic("implement me")
}

func (r *ApartmentResource) Update(id string, jsonData []byte) ([]byte, error) {
	panic("implement me")
}

func (r *ApartmentResource) Delete(id string) error {
	panic("implement me")
}
