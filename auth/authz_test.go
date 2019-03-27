package auth

import (
	"testing"
)

func TestOnlyReadForOneRole(t *testing.T) {
	// Arrange
	authorizer := NewAuthzService()
	authorizer.AddPermission("client", "user", Read)

	// Act & Assert
	assert(t, authorizer.Allowed("client", "user", Read))
	assert(t, !authorizer.Allowed("client", "user", Create))
	assert(t, !authorizer.Allowed("client", "user", Delete))
	assert(t, !authorizer.Allowed("client", "user", Update))
}

func TestCreateReadTwoRoles(t *testing.T) {
	// Arrange
	authorizer := NewAuthzService()
	authorizer.AddPermission("client", "user", Read)
	authorizer.AddPermission("realtor", "user", Create, Read, Update)

	// Act & Assert
	assert(t, authorizer.Allowed("client", "user", Read))
	assert(t, !authorizer.Allowed("client", "user", Create))
	assert(t, !authorizer.Allowed("client", "user", Delete))
	assert(t, !authorizer.Allowed("client", "user", Update))

	assert(t, authorizer.Allowed("realtor", "user", Create))
	assert(t, authorizer.Allowed("realtor", "user", Read))
	assert(t, authorizer.Allowed("realtor", "user", Update))
	assert(t, !authorizer.Allowed("realtor", "user", Delete))
}

func TestTwoResources(t *testing.T) {
	// Arrange
	authorizer := NewAuthzService()
	authorizer.AddPermission("client", "user", Read)
	authorizer.AddPermission("client", "loser", Create)

	// Act & Assert
	assert(t, authorizer.Allowed("client", "user", Read))
	assert(t, !authorizer.Allowed("client", "user", Create))
	assert(t, !authorizer.Allowed("client", "user", Delete))
	assert(t, !authorizer.Allowed("client", "user", Update))
	assert(t, authorizer.Allowed("client", "loser", Create))
	assert(t, !authorizer.Allowed("client", "loser", Update))
	assert(t, !authorizer.Allowed("client", "loser", Read))
	assert(t, !authorizer.Allowed("client", "loser", Delete))
}

func assert(t *testing.T, expr bool) {
	t.Helper()

	if !expr {
		t.Error("Expression is false")
	}
}
