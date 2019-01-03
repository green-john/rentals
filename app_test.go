// +build integration

package rentals

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
)

func TestCreateUser(t *testing.T) {
	// Arrange
	var wg sync.WaitGroup
	const port = 8087
	app, err := NewApp(port)
	ok(t, err)

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Serve()
	}()

	payload := createPayload(`{"username":"john", "password": "secret", "role": "client"}`)
	serverUrl := app.getServerURL()
	url := fmt.Sprintf("http://%s/users", serverUrl)

	t.Run("create and get user", func(t *testing.T) {
		defer app.dropDB()
		// Act
		res, err := http.Post(url, "application/json", payload)
		ok(t, err)

		// Assert
		assert(t, res.StatusCode == http.StatusCreated, fmt.Sprintf("Expected 201 got %d", res.StatusCode))
		rawContent, err := ioutil.ReadAll(res.Body)
		ok(t, err)

		var userResponse struct {
			Username string `json:"username"`
			PwdHash  string `json:"password_hash"`
			Role     string `json:"role"`
		}

		err = json.Unmarshal(rawContent, &userResponse)

		assert(t, userResponse.Username == "john",
			fmt.Sprintf("Expected name john, got %s", userResponse.Username))
		assert(t, userResponse.PwdHash == "hash[secret]",
			fmt.Sprintf("Expected pass hash[secret], got %s", userResponse.PwdHash))
		assert(t, userResponse.Role == "client",
			fmt.Sprintf("Expected role client, got %s", userResponse.Role))
	})
}

func createPayload(s string) *bytes.Buffer {
	return bytes.NewBuffer([]byte(s))
}
