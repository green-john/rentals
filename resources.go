package rentals

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
)

type NotFoundError string

func (e NotFoundError) Error() string {
	return string(e)
}

type UserResource struct {
	Db *gorm.DB
}

func (ur *UserResource) Name() string {
	return "users"
}

func (ur *UserResource) Create(jsonData []byte) ([]byte, error) {
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

	ur.Db.Create(&user)

	rawJson, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	return rawJson, nil
}

func (ur *UserResource) All() ([]byte, error) {
	var users []User
	ur.Db.Find(&users)

	return json.Marshal(users)
}

func (ur *UserResource) Read(id string) ([]byte, error) {
	user, err := getUserIdStr(id, ur.Db)
	if err != nil {
		return nil, err
	}
	return json.Marshal(user)
}

func (ur *UserResource) Update(id string, jsonData []byte) ([]byte, error) {
	user, err := getUserIdStr(id, ur.Db)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &user)
	if err != nil {
		return nil, err
	}

	// Save to DB
	ur.Db.Save(&user)
	return json.Marshal(user)
}

func (ur *UserResource) Delete(id string) error {
	user, err := getUserIdStr(id, ur.Db)
	if err != nil {
		return err
	}

	ur.Db.Delete(&user)
	return nil
}

type ApartmentResource struct {
	Db *gorm.DB
}

func (ar *ApartmentResource) Name() string {
	return "apartments"
}

func (ar *ApartmentResource) Create(jsonData []byte) ([]byte, error) {
	newApartment, err := ar.createApartment(jsonData)
	if err != nil {
		return nil, err
	}

	ar.Db.Create(&newApartment)

	rawJson, err := json.Marshal(newApartment)
	if err != nil {
		return nil, err
	}

	return rawJson, nil
}

func (ar *ApartmentResource) Read(id string) ([]byte, error) {
	apartment, err := getApartment(id, ar.Db)
	if err != nil {
		return nil, err
	}

	return json.Marshal(apartment)
}

func (ar *ApartmentResource) All() ([]byte, error) {
	var apartments []Apartment
	ar.Db.Find(&apartments)

	return json.Marshal(apartments)
}

func (ar *ApartmentResource) Update(id string, jsonData []byte) ([]byte, error) {
	apartment, err := getApartment(id, ar.Db)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &apartment)
	if err != nil {
		return nil, err
	}

	// Save to DB
	ar.Db.Save(&apartment)

	return json.Marshal(apartment)
}

func (ar *ApartmentResource) Delete(id string) error {
	apartment, err := getApartment(id, ar.Db)
	if err != nil {
		return err
	}

	ar.Db.Delete(&apartment)
	return nil
}

func getApartment(id string, db *gorm.DB) (*Apartment, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var apartment Apartment
	db.First(&apartment, idInt)

	if apartment.ID != uid(idInt) {
		return nil, NotFoundError(fmt.Sprintf("apartment %d not found", idInt))
	}

	return &apartment, nil
}

func (ar *ApartmentResource) createApartment(jsonData []byte) (*Apartment, error) {
	var newApartment Apartment

	err := json.Unmarshal(jsonData, &newApartment)
	if err != nil {
		return nil, fmt.Errorf("[createApartment] error calling json.Unmarshall(): %v", err)
	}

	realtor := getUser(newApartment.RealtorId, ar.Db)
	if realtor == nil {
		return nil, fmt.Errorf("realtor not found (id=%d)", newApartment.RealtorId)
	}

	err = newApartment.Validate()
	if err != nil {
		return nil, err
	}

	return &newApartment, nil
}

func getUserIdStr(id string, db *gorm.DB) (*User, error) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	usr := getUser(uint(intId), db)
	if usr == nil {
		return nil, NotFoundError(fmt.Sprintf("user %d not found", intId))
	}

	return usr, nil
}

func getUser(userId uint, db *gorm.DB) *User {
	var user User
	db.First(&user, userId)

	if user.ID != uid(userId) {
		return nil
	}

	return &user
}
