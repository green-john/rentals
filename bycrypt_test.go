package rentals

import (
	"fmt"
	"testing"
)

func TestEncryptPassword(t *testing.T) {
	// Arrange
	clearPass := "password"

	// Act
	encrypted, err := EncryptPassword(clearPass)
	ok(t, err)

	// Assert
	err = CheckPassword(encrypted, "password")
	assert(t, err == nil, fmt.Sprintf("Unexpected error %v", err))
}
