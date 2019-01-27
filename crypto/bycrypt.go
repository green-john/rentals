package crypto

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hashPassword), nil
}

func CheckPassword(hashPassword, password string) error {
	hashBytes, err := hex.DecodeString(hashPassword)
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword(hashBytes, []byte(password))
}
