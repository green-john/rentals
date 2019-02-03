package crypto

import (
	"fmt"
	"rentals/tst"
	"testing"
)

func TestEncryptPassword(t *testing.T) {
	// Arrange
	clearPass := "password"

	// Act
	encrypted, err := EncryptPassword(clearPass)
	tst.Ok(t, err)

	// True
	err = CheckPassword(encrypted, "password")
	tst.True(t, err == nil, fmt.Sprintf("Unexpected error %v", err))
}
