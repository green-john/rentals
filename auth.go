package rentals

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

// Elements related to authentication and authorization.

// Error thrown when a login fails
var LoginError = errors.New("incorrect username/password")

// Authenticator is the interface that should be implemented when
// designing an auth scheme that depends on a stateful bearer token
// Note: It is NOT safe to use this for stateless authentications schemes
// such as Jose/JWT/Macaroon.
type Authenticator interface {
	// Login tries to login a user given its username and password.
	// logging in a user entails checking whether the info is correct
	// and in case it is, generate a token that can be used
	// for future requests. Users should include this token in
	// their requests
	Login(username, password string) (string, error)

	// Verify checks whether or not the given token is valid.
	// If it is, it returns the user associated to such token.
	// Otherwise, returns nil.
	Verify(token string) *User
}

// Implementation of a Authenticator using a relational database
type dbAuthenticator struct {
	Db *gorm.DB
}

func (a *dbAuthenticator) Login(username, password string) (string, error) {
	var user User
	a.Db.Where("username = ?", username).First(&user)

	// Username was not found as we don't allow empty passwords
	if user.PasswordHash == "" {
		return "", LoginError
	}

	if CheckPassword(user.PasswordHash, password) != nil {
		return "", LoginError
	}

	// Check if there is an existing session already.
	existingSession := a.findExistingUserSession(user)
	if existingSession != nil {
		return existingSession.Token, nil
	}

	// Otherwise, create a new token and session and save it to the db
	token := generateToken()
	session := UserSession{
		Token:  token,
		UserID: user.ID,
		User:   user,
	}
	a.Db.Create(&session)

	return token, nil
}

func (a *dbAuthenticator) findExistingUserSession(user User) *UserSession {
	var session UserSession

	a.Db.Where("user_id = ?", user.ID).First(&session)

	// Session not found
	if session.Token == "" {
		return nil
	}

	return &session
}

func generateToken() string {
	const tokenLength = 24
	ret := make([]byte, tokenLength)
	_, err := rand.Read(ret)

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%X", ret)
}

func (a *dbAuthenticator) Verify(token string) *User {
	var userSession UserSession
	a.Db.Where("token = ?", token).First(&userSession)

	if userSession.Token != token {
		return nil
	}

	var user User
	a.Db.Model(&userSession).Related(&user)

	return &user
}

// Creates a new database authenticator
func NewDbAuthenticator(db *gorm.DB) *dbAuthenticator {
	return &dbAuthenticator{Db: db}
}
