package rentals

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
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

func (t *UserResource) All() ([]byte, error) {
	var users []User
	t.Db.Find(&users)

	return json.Marshal(users)
}

func (t *UserResource) Read(id string) ([]byte, error) {
	var user User
	t.Db.First(&user, id)

	return json.Marshal(user)
}

func (t *UserResource) Update(id string, jsonData []byte) ([]byte, error) {
	panic("implement me")
}

func (t *UserResource) Delete(id string) error {
	panic("implement me")
}

func (r *ApartmentResource) getUser(userId uint) *User {
	var user User
	r.Db.First(&user, userId)

	if user.ID != userId {
		return nil
	}

	return &user
}

type ApartmentResource struct {
	Db *gorm.DB
}

func (r *ApartmentResource) Name() string {
	return "apartments"
}

func (r *ApartmentResource) Create(jsonData []byte) ([]byte, error) {
	newApartment, err := r.createApartment(jsonData)
	if err != nil {
		return nil, err
	}

	r.Db.Create(&newApartment)

	rawJson, err := json.Marshal(newApartment)
	if err != nil {
		return nil, err
	}

	return rawJson, nil
}

func (r *ApartmentResource) Read(id string) ([]byte, error) {
	apartment, err := r.getApartment(id)
	if err != nil {
		return nil, err
	}

	return json.Marshal(apartment)
}

func (r *ApartmentResource) All() ([]byte, error) {
	var apartments []Apartment
	r.Db.Find(&apartments)

	return json.Marshal(apartments)
}

func (r *ApartmentResource) Update(id string, jsonData []byte) ([]byte, error) {
	apartment, err := r.getApartment(id)
	if err != nil {
		return nil, err
	}

	//cleanedJson, err := filterUpdateFields(jsonData)
	//if err != nil {
	//	return nil, err
	//}

	err = json.Unmarshal(jsonData, &apartment)
	if err != nil {
		return nil, err
	}

	// Save to DB
	r.Db.Save(&apartment)

	rawJson, err := json.Marshal(apartment)
	if err != nil {
		return nil, err
	}

	return rawJson, nil
}

func (r *ApartmentResource) Delete(id string) error {
	apartment, err := r.getApartment(id)
	if err != nil {
		return err
	}

	r.Db.Delete(&apartment)
	return nil
}

func (r *ApartmentResource) getApartment(id string) (*Apartment, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var apartment Apartment
	r.Db.First(&apartment, idInt)

	if apartment.ID != uid(idInt) {
		return nil, errors.New("apartment not found")
	}

	return &apartment, nil
}

func (r *ApartmentResource) createApartment(jsonData []byte) (*Apartment, error) {
	var newApartment Apartment

	err := json.Unmarshal(jsonData, &newApartment)
	if err != nil {
		return nil, fmt.Errorf("[createApartment] error calling json.Unmarshall(): %v", err)
	}

	realtor := r.getUser(newApartment.RealtorId)
	if realtor == nil {
		return nil, fmt.Errorf("realtor not found (id=%d)", newApartment.RealtorId)
	}

	err = newApartment.Validate()
	if err != nil {
		return nil, err
	}

	return &newApartment, nil
}
