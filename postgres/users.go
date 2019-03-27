package postgres

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"rentals"
	"rentals/crypto"
	"strconv"
)


type dbUserService struct {
	Db *gorm.DB
}

func (s *dbUserService) Create(input rentals.UserCreateInput) (*rentals.UserCreateOutput, error) {
	user, err := createUser(input.Username, input.Password, input.Role, s.Db)
	if err != nil {
		return nil, err
	}
	return &rentals.UserCreateOutput{User: *user}, nil
}

func (s *dbUserService) All(rentals.UserAllInput) (*rentals.UserAllOutput, error) {
	var users []rentals.User
	if err := s.Db.Find(&users).Error; err != nil {
		return nil, err
	}

	return &rentals.UserAllOutput{Users: users}, nil
}

func (s *dbUserService) Read(input rentals.UserReadInput) (*rentals.UserReadOutput, error) {
	user, err := getUser(input.Id, s.Db)
	if err != nil {
		return nil, err
	}

	return &rentals.UserReadOutput{User: *user}, nil
}

func (s *dbUserService) Update(input rentals.UserUpdateInput) (*rentals.UserUpdateOutput, error) {
	user, err := getUser(input.Id, s.Db)
	if err != nil {
		return nil, err
	}

	if input.Password != "" {
		user.PasswordHash, err = crypto.EncryptPassword(input.Password)
		if err != nil {
			return nil, fmt.Errorf("[dbUserService.Update] error encrypting password %v", err)
		}
	}

	if input.Role != "" && validRole(input.Role) {
		user.Role = input.Role
	}

	// Save to DB
	if err := s.Db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("[dbUserService.Update] error updating %v", err)
	}

	return &rentals.UserUpdateOutput{User: *user}, nil
}

func (s *dbUserService) Delete(input rentals.UserDeleteInput) (*rentals.UserDeleteOutput, error) {
	user, err := getUser(input.Id, s.Db)
	if err != nil {
		return nil, err
	}

	if err := s.Db.Delete(&user).Error; err != nil {
		return nil, err
	}

	return &rentals.UserDeleteOutput{}, nil
}

func NewDbUserService(db *gorm.DB) *dbUserService {
	return &dbUserService{Db: db}
}

func getUser(id string, db *gorm.DB) (*rentals.User, error) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var user rentals.User

	if err = db.First(&user, intId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, rentals.NotFoundError
		}
		return nil, err
	}

	return &user, nil
}

func createUser(username, password, role string, db *gorm.DB) (*rentals.User, error) {
	if !validRole(role) {
		return nil, errors.New(
			fmt.Sprintf("error creating user. Unknown role %s", role))
	}

	pwdHash, err := crypto.EncryptPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error encrypting password %v", err)
	}

	user := rentals.User{
		Username:     username,
		PasswordHash: pwdHash,
		Role:         role,
	}

	if err = db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("error creating user %v", err)
	}

	return &user, nil
}

func validRole(role string) bool {
	return contains([]string{"admin", "realtor", "client"}, role)
}

func contains(a []string, b string) bool {
	for _, elt := range a {
		if elt == b {
			return true
		}
	}

	return false
}
