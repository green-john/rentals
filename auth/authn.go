package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"rentals"
	"rentals/crypto"
)

// Elements related to authentication and authorization.

// Error thrown when a login fails
var LoginError = errors.New("incorrect username/password")

// AuthnService is the interface that should be implemented when
// designing an auth scheme that depends on a stateful bearer token
// Note: It is NOT safe to use this for stateless authentications schemes
// such as Jose/JWT/Macaroon.
type AuthnService interface {
	// Login tries to login a user given its username and password.
	// logging in a user entails checking whether the info is correct
	// and in case it is, generate a token that can be used
	// for future requests. Users should include this token in
	// their requests
	Login(username, password string) (string, error)

	// Verify checks whether or not the given token is valid.
	// If it is, it returns the user associated to such token.
	// Otherwise, returns nil.
	Verify(token string) *rentals.User
}

// Implementation of a AuthnService using a relational database
type dbAuthnService struct {
	Db *gorm.DB
}

func (a *dbAuthnService) Login(username, password string) (string, error) {
	var user rentals.User
	a.Db.Where("username = ?", username).First(&user)

	// Username was not found as we don't allow empty passwords
	if user.PasswordHash == "" {
		return "", LoginError
	}

	if crypto.CheckPassword(user.PasswordHash, password) != nil {
		return "", LoginError
	}

	// Check if there is an existing session already.
	existingSession := a.findExistingUserSession(user)
	if existingSession != nil {
		return existingSession.Token, nil
	}

	// Otherwise, create a new token and session and save it to the Db
	token := generateToken()
	session := rentals.UserSession{
		Token:  token,
		UserID: uint(user.ID),
		User:   user,
	}
	a.Db.Create(&session)

	return token, nil
}

func (a *dbAuthnService) findExistingUserSession(user rentals.User) *rentals.UserSession {
	var session rentals.UserSession

	a.Db.Where("user_id = ?", user.ID).First(&session)

	// Session not found
	if session.Token == "" {
		return nil
	}

	return &session
}

// Generates a random by drawing a number of bytes from
// crypto.Rand
func generateToken() string {
	const tokenLength = 24
	ret := make([]byte, tokenLength)
	_, err := rand.Read(ret)

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%X", ret)
}

func (a *dbAuthnService) Verify(token string) *rentals.User {
	var userSession rentals.UserSession
	a.Db.Where("token = ?", token).First(&userSession)

	if userSession.Token != token {
		return nil
	}

	var user rentals.User
	a.Db.Model(&userSession).Related(&user)

	return &user
}

// Creates a new database authenticator
func NewDbAuthnService(db *gorm.DB) *dbAuthnService {
	return &dbAuthnService{Db: db}
}
