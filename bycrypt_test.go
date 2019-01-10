package rentals

import (
	"fmt"
	"testing"
	"tournaments/tst"
)

func TestEncryptPassword(t *testing.T) {
	// Arrange
	clearPass := "password"

	// Act
	encrypted, err := EncryptPassword(clearPass)
	tst.Ok(t, err)

	// Assert
	err = CheckPassword(encrypted, "password")
	tst.Assert(t, err == nil, fmt.Sprintf("Unexpected error %v", err))
}
