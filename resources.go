package rentals

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

type UserResource struct {
	DB *gorm.DB
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

	pwdHash := fmt.Sprintf("hash[%s]", newUserSchema.Password)

	user := User{
		Username:     newUserSchema.Username,
		PasswordHash: pwdHash,
		Role:         newUserSchema.Role,
	}

	t.DB.Create(&user)

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
